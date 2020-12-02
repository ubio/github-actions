package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/go-github/v30/github"
	"golang.org/x/oauth2"
)

var (
	err    error
	cfg    config
	client *github.Client
	ctx    context.Context
)

type config struct {
	Token   string `env:"INPUT_TOKEN,required"`
	Owner   string `env:"INPUT_OWNER,required"`
	Repo    string `env:"INPUT_REPOSITORY,required"`
	Message string `env:"INPUT_MESSAGE,required"`
	Files   string `env:"INPUT_FILES,required"`

	// PR Vars
	Title               string   `env:"INPUT_TITLE,required"`
	Head                string   `env:"INPUT_HEAD,required"`
	Base                string   `env:"INPUT_BASE,required"`
	Body                string   `env:"INPUT_BODY" envDefault:""`
	MaintainerCanModify bool     `env:"INPUT_MAINTAINER_CAN_MODIFY" envDefault:"true"`
	Draft               bool     `env:"INPUT_DRAFT" envDefault:"false"`
	Merge               bool     `env:"INPUT_MERGE" envDefault:"false"`
	Labels              []string `env:"INPUT_LABELS" envDefault:""`
}

func init() {
	if err = env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	ctx = context.Background()
	client = github.NewClient(
		oauth2.NewClient(
			ctx, oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: cfg.Token,
				},
			),
		),
	)
}

func main() {

	if cfg.Files == "" {
		log.Fatal("No files to commit")
	}

	ref, err := getRef()
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %s\n", err)
	}
	if ref == nil {
		log.Fatalf("No error returned but the github reference is `nil`")
	}

	tree, err := getTree(ref)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}

	if err = pushCommit(ref, tree); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}

	pr, err := createPR()
	if err != nil {
		log.Fatalf("Error while creating the pull request: %s", err)
	}
	log.Println("PR created:", pr.GetHTMLURL())

	if err := addLabelsToPR(cfg.Labels, pr); err != nil {
		log.Fatalf("Error while adding labels to the pull request: %s", err)
	}

	fmt.Println(fmt.Sprintf(`::set-output name=pr::%s`, pr.GetHTMLURL()))

	if cfg.Merge {
		if err := mergePullRequest(pr); err != nil {
			fmt.Println("::set-output name=merged::false")
			log.Fatalf("error merging pull request, %s", err.Error())
		}
		fmt.Println("::set-output name=merged::true")
		log.Println("successfully merged pull request")
	}
}

func addLabelsToPR(labels []string, pr *github.PullRequest) error {
	_, _, err := client.Issues.AddLabelsToIssue(ctx, cfg.Owner, cfg.Repo, pr.GetNumber(), labels)
	return err
}

// createPR builds and creates the PR on github
func createPR() (*github.PullRequest, error) {

	pr, _, err := client.PullRequests.Create(
		ctx,
		cfg.Owner,
		cfg.Repo,
		&github.NewPullRequest{
			Title:               &cfg.Title,
			Head:                &cfg.Head,
			Base:                &cfg.Base,
			Body:                &cfg.Body,
			MaintainerCanModify: &cfg.MaintainerCanModify,
			Draft:               &cfg.Draft,
		},
	)

	return pr, err
}

// getRef tries to fetch the reference and if it can't be found, creates one
func getRef() (ref *github.Reference, err error) {
	ref, _, err = client.Git.GetRef(
		ctx,
		cfg.Owner,
		cfg.Repo,
		fmt.Sprintf("refs/heads/%s", cfg.Head),
	)
	if err == nil {
		return ref, nil
	}

	if cfg.Head == cfg.Base {
		return nil, errors.New("`base` is the same as `head`")
	}

	var baseRef *github.Reference
	baseRef, _, err = client.Git.GetRef(
		ctx,
		cfg.Owner,
		cfg.Repo,
		fmt.Sprintf("refs/heads/%s", cfg.Base),
	)
	if err != nil {
		return nil, err
	}

	newRef := &github.Reference{
		Ref: github.String("refs/heads/" + cfg.Head),
		Object: &github.GitObject{
			SHA: baseRef.Object.SHA,
		},
	}
	ref, _, err = client.Git.CreateRef(ctx, cfg.Owner, cfg.Repo, newRef)

	return ref, err
}

// pushCommit creates the commit in the given reference using the given tree
func pushCommit(ref *github.Reference, tree *github.Tree) (err error) {

	// get the parent commit to attach the commit to
	parent, _, err := client.Repositories.GetCommit(ctx, cfg.Owner, cfg.Repo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	parent.Commit.SHA = parent.SHA

	// create the commit using the tree
	newCommit, _, err := client.Git.CreateCommit(
		ctx,
		cfg.Owner,
		cfg.Repo,
		&github.Commit{
			Message: &cfg.Message,
			Tree:    tree,
			Parents: []*github.Commit{
				parent.Commit,
			},
		},
	)
	if err != nil {
		return err
	}

	// attach the commit to the master branch
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, cfg.Owner, cfg.Repo, ref, false)
	return err
}

func getTree(ref *github.Reference) (tree *github.Tree, err error) {

	entries := []*github.TreeEntry{}

	for _, file := range strings.Split(cfg.Files, ",") {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		entries = append(
			entries,
			&github.TreeEntry{
				Path:    github.String(file),
				Type:    github.String("blob"),
				Content: github.String(string(content)),
				Mode:    github.String("100644"),
			},
		)
	}

	tree, _, err = client.Git.CreateTree(
		ctx,
		cfg.Owner,
		cfg.Repo,
		*ref.Object.SHA,
		entries,
	)
	return tree, err
}

func mergePullRequest(pr *github.PullRequest) error {
	if err = awaitMergeableState(pr); err != nil {
		return err
	}
	msg := "Auto merging pull request"
	state, resp, err := client.PullRequests.Merge(ctx, cfg.Owner, cfg.Repo, pr.GetNumber(), msg, nil)
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("nil response received")
	}
	if !*state.Merged {
		fmt.Errorf("failed to merge pull request")
	}
	return nil
}

func awaitMergeableState(pr *github.PullRequest) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 30 * time.Second
	bo.InitialInterval = 3 * time.Second
	bo.MaxInterval = 10 * time.Second

	ticker := backoff.NewTicker(bo)

	for range ticker.C {

		state, resp, err := client.PullRequests.Get(ctx, cfg.Owner, cfg.Repo, pr.GetNumber())
		if err != nil {
			log.Printf("An error occurred talking to the GitHub API, %s\n", err.Error())
			continue
		}

		if resp == nil {
			log.Printf("A nil response was returned by the GitHub API")
			continue
		}

		// we could be a bit cleverer here and break
		// early on various GetMergeableState values
		if !state.GetMergeable() {
			log.Println("pull request is not yet mergeable")
			continue
		}

		ticker.Stop()
		return nil
	}

	return fmt.Errorf("timed out waiting for PR to be mergeable")
}
