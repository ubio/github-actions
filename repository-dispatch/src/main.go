package main

import (
	"context"
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	cfg config
)

type config struct {
	Token   string `env:"INPUT_TOKEN"`
	Repo    string `env:"INPUT_REPOSITORY"`
	Event   string `env:"INPUT_EVENT"`
	Payload string `env:"INPUT_PAYLOAD"`
}

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)
}

func main() {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(repos)
}
