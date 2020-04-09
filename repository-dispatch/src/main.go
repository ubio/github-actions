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
	cfg config
)

type config struct {
	Token   string `env:"INPUT_TOKEN"`
	Owner   string `env:"INPUT_OWNER"`
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

	client, ctx := buildClient()

	_, resp, err := client.Repositories.Dispatch(ctx, cfg.Owner, cfg.Repo, buildDispatchRequestOptions())
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

func buildDispatchRequestOptions() github.DispatchRequestOptions {
	msg := json.RawMessage([]byte(cfg.Payload))

	return github.DispatchRequestOptions{
		EventType:     cfg.Event,
		ClientPayload: &msg,
	}
}
