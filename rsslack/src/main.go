package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
)

func main() {

	token, err := env("INPUT_SLACK_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	channel, err := env("INPUT_CHANNEL")
	if err != nil {
		log.Fatal(err)
	}

	message, err := env("INPUT_MESSAGE")
	if err != nil {
		log.Fatal(err)
	}

	api := slack.New(token)
	channelID, timestamp, err := api.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Message sent to channel %s at %s", channelID, timestamp)
}

func env(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("env key %s does not exist", key)
	}
	if len(value) == 0 {
		return "", fmt.Errorf("env key %s is empty", key)
	}
	return value, nil
}
