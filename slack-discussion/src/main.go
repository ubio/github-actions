package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/kelseyhightower/envconfig"
)

var (
	vars EnvVars
	err  error
)

// EnvVars passed by GH actions
type EnvVars struct {
	Token     string `envconfig:"INPUT_SLACK_TOKEN" required:"true"`
	Channel   string `envconfig:"INPUT_CHANNEL" required:"true"`
	EventPath string `envconfig:"GITHUB_EVENT_PATH" required:"true"`
}

type DiscussionEvent struct {
	Discussion *Discussion `json:"discussion,omitempty"`
}

type Discussion struct {
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	CreatedAt string    `json:"created_at,omitempty"`
	HTMLURL   string    `json:"html_url,omitempty"`
	Category  *Category `json:"category,omitempty"`
	User      *User     `json:"user,omitempty"`
}

type Category struct {
	Emoji       string `json:"emoji,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type User struct {
	AvatarURL string `json:"avatar_url,omitempty"`
	HTMLURL   string `json:"html_url,omitempty"`
}

func main() {

	err = envconfig.Process("", &vars)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(vars.EventPath)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	discussion := DiscussionEvent{}
	err = json.Unmarshal(content, &discussion)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(discussion)

	// msg := fmt.Sprintf(template)

	// api := slack.New(vars.Token)
	// channelID, timestamp, err := api.PostMessage(vars.Channel, slack.MsgOptionText(msg, false))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("Message sent to channel %s at %s", channelID, timestamp)
}
