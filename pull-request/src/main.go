package main

import (
	"context"
	"log"

	"github.com/caarlos0/env"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-github/v30/github"
	"golang.org/x/oauth2"
)

var (
	cfg config
)

type config struct {
	Token   string `env:"INPUT_TOKEN,required"`
	Owner   string `env:"INPUT_OWNER,required"`
	Repo    string `env:"INPUT_REPOSITORY,required"`
	Message string `env:"INPUT_MESSAGE,required"`
	SHA     string `env:"INPUT_GIT_SHA,required"`

	// PR Vars
	Title               string `env:"INPUT_TITLE,required"`
	Head                string `env:"INPUT_HEAD,required"`
	Base                string `env:"INPUT_BASE,required"`
	Body                string `env:"INPUT_BODY" envDefault:""`
	MaintainerCanModify bool   `env:"INPUT_MAINTAINER_CAN_MODIFY" envDefault:"true"`
	Draft               bool   `env:"INPUT_DRAFT" envDefault:"false"`
}

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
}

func main() {

	client, ctx := buildClient()

	tree, _, err := client.Git.GetTree(ctx, cfg.Owner, cfg.Repo, cfg.SHA, false)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(tree)
	c, _, err := client.Git.CreateCommit(ctx, cfg.Owner, cfg.Repo, buildCommit(tree))
	spew.Dump(c)
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = client.PullRequests.Create(ctx, cfg.Owner, cfg.Repo, buildPullRequest())
	if err != nil {
		log.Fatal(err)
	}
}

func buildClient() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: cfg.Token,
		},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), ctx
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

func buildCommit(tree *github.Tree) *github.Commit {
	return &github.Commit{
		Message: &cfg.Message,
		Tree:    tree,
	}
}
