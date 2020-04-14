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
	SHA     string `env:"INPUT_GIT_SHA,required"`
	Files   string `env:"INPUT_FILES"`

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
	initClient()
}

func initClient() {
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

	// @todo add check
	if cfg.Files == "" {
		// skip branch + commit
	}

	ref, err := getRef()
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %s\n", err)
	}
	if ref == nil {
		log.Fatalf("No error where returned but the reference is nil")
	}

	tree, err := getTree(ref)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := pushCommit(ref, tree); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}

	if err := createPR(); err != nil {
		log.Fatalf("Error while creating the pull request: %s", err)
	}
}

func buildPullRequest() *github.NewPullRequest {
	return &github.NewPullRequest{
		Title:               github.String(cfg.Title),
		Head:                github.String(cfg.Head),
		Base:                github.String(cfg.Base),
		Body:                github.String(cfg.Body),
		MaintainerCanModify: github.Bool(cfg.MaintainerCanModify),
		Draft:               github.Bool(cfg.Draft),
	}
}

func createPR() (err error) {

	newPR := &github.NewPullRequest{
		Title:               &cfg.Title,
		Head:                &cfg.Head,
		Base:                &cfg.Base,
		Body:                &cfg.Body,
		MaintainerCanModify: &cfg.MaintainerCanModify,
	}

	pr, _, err := client.PullRequests.Create(ctx, cfg.Owner, cfg.Repo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}

func buildCommit(tree *github.Tree) *github.Commit {
	return &github.Commit{
		Message: &cfg.Message,
		Tree:    tree,
	}
}

// createCommit creates the commit in the given reference using the given tree.
func pushCommit(ref *github.Reference, tree *github.Tree) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, cfg.Owner, cfg.Repo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &date}
	commit := &github.Commit{Author: author, Message: &cfg.Message, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, cfg.Owner, cfg.Repo, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, cfg.Owner, cfg.Repo, ref, false)
	return err
}

func getRef() (ref *github.Reference, err error) {
	if ref, _, err = client.Git.GetRef(ctx, cfg.Owner, cfg.Repo, "refs/heads/"+cfg.Head); err == nil {
		return ref, nil
	}

	if cfg.Head == cfg.Base {
		return nil, errors.New("`base` is the same as `head`")
	}

	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, cfg.Owner, cfg.Repo, "refs/heads/"+cfg.Base); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + cfg.Head), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = client.Git.CreateRef(ctx, cfg.Owner, cfg.Repo, newRef)
	return ref, err
}

func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("No files supplied")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = ioutil.ReadFile(localFile)
	return targetName, b, err
}

func getTree(ref *github.Reference) (tree *github.Tree, err error) {

	entries := []*github.TreeEntry{}

	for _, fileArg := range strings.Split(cfg.Files, ",") {
		file, content, err := getFileContent(fileArg)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(content)), Mode: github.String("100644")})
	}

	tree, _, err = client.Git.CreateTree(ctx, cfg.Owner, cfg.Repo, *ref.Object.SHA, entries)
	return tree, err
}
