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
	token   string      `env:"INPUT_TOKEN"`
	repo    string      `env:"INPUT_REPOSITORY"`
	event   string      `env:"INPUT_EVENT"`
	payload interface{} `env:"INPUT_PAYLOAD"`
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
