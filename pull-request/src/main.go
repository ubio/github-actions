package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/caarlos0/env"
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
	Title               string `env:"INPUT_TITLE,required"`
	Head                string `env:"INPUT_HEAD,required"`
	Base                string `env:"INPUT_BASE,required"`
	Body                string `env:"INPUT_BODY" envDefault:""`
	MaintainerCanModify bool   `env:"INPUT_MAINTAINER_CAN_MODIFY" envDefault:"true"`
	Draft               bool   `env:"INPUT_DRAFT" envDefault:"false"`
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

	fmt.Println(fmt.Sprintf(`::set-output name=pr::%s`, pr.GetHTMLURL()))
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
