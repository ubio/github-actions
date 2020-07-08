package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

func main() {

	api := slack.New(os.Getenv("INPUT_SLACK_TOKEN"))
	channelID, timestamp, err := api.PostMessage(os.Getenv("INPUT_CHANNEL"), slack.MsgOptionText(os.Getenv("INPUT_MESSAGE"), false))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Message sent to channel %s at %s", channelID, timestamp)
}
