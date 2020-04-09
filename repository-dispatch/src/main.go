package main

import (
	"context"
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
}

func main() {

	client := buildClient()

	resp, err := client.Repositories.Dispatch(ctx, "", cfg.Repo, buildDispatchRequestOptions())
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}

func buildClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: cfg.Token,
		},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func buildDispatchRequestOptions() github.DispatchRequestOptions {
	msg := []byte(cfg.Payload)

	return github.DispatchRequestOptions{
		EventType:     cfg.Event,
		ClientPayload: &msg,
	}
}
