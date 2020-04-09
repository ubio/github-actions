package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-github/github"
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

	client := github.NewClient(nil)

	spew.Dump(client)
}
