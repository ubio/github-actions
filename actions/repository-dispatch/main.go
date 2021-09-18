package main

import (
	"context"
	"encoding/json"
	"log"

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
	Event   string `env:"INPUT_EVENT,required"`
	Payload string `env:"INPUT_PAYLOAD"`
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
	_, _, err = client.Repositories.Dispatch(ctx, cfg.Owner, cfg.Repo, buildDispatchRequestOptions())
	if err != nil {
		log.Fatal(err)
	}
}

func buildDispatchRequestOptions() github.DispatchRequestOptions {
	msg := json.RawMessage([]byte(cfg.Payload))

	return github.DispatchRequestOptions{
		EventType:     cfg.Event,
		ClientPayload: &msg,
	}
}
