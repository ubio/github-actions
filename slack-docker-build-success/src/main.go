package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

const template = "*%s* has been built :package:\n<%s/%s/%s:%s|%s/%s:%s>"

var (
	vars EnvVars
	err  error
)

// EnvVars passed by GH actions
type EnvVars struct {
	Token   string `envconfig:"INPUT_SLACK_TOKEN" required:"true"`
	Channel string `envconfig:"INPUT_CHANNEL" required:"true"`
	Name    string `envconfig:"INPUT_NAME" required:"true"`

	Registry  string `envconfig:"INPUT_REGISTRY" required:"true"`
	Namespace string `envconfig:"INPUT_NAMESPACE" required:"true"`
	Image     string `envconfig:"INPUT_IMAGE" required:"true"`
	Tag       string `envconfig:"INPUT_TAG" required:"true"`
}

func main() {

	err = envconfig.Process("", &vars)
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf(template, vars.Name, vars.Registry, vars.Namespace, vars.Image, vars.Tag, vars.Namespace, vars.Image, vars.Tag)

	api := slack.New(vars.Token)
	channelID, timestamp, err := api.PostMessage(vars.Channel, slack.MsgOptionText(msg, false))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Message sent to channel %s at %s", channelID, timestamp)
}
