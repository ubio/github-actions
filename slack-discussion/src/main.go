package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/eritikass/githubmarkdownconvertergo"
	"github.com/kelseyhightower/envconfig"
	"github.com/slack-go/slack"
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
	Login     string `json:"login,omitempty"`
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

	event := DiscussionEvent{}
	err = json.Unmarshal(content, &event)
	if err != nil {
		log.Fatal(err)
	}

	api := slack.New(vars.Token)
	channelID, timestamp, err := api.PostMessage(
		vars.Channel,
		slack.MsgOptionText(buildSlackBlock(event.Discussion), true),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func buildSlackBlock(d *Discussion) slack.Message {

	dividerSection := slack.NewDividerBlock()

	// @TODO: get the squad name from the event
	squadName := "Proxies Squad"

	date, err := time.Parse(
		time.RFC3339,
		d.CreatedAt,
	)
	if err != nil {
		log.Fatal(err)
	}

	day := date.Format("Mon, Jan 2")
	tod := date.Format("03:04")

	contextHeaderText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("%s %s - *%s*\n%s @ %s by <%s|%s>", d.Category.Emoji, squadName, d.Category.Name, day, tod, d.User.HTMLURL, d.User.Login), false, false)

	buttonText := slack.NewTextBlockObject("plaintext", "View Discussion", false, false)
	button := slack.NewButtonBlockElement("discussion-link", d.HTMLURL, buttonText)
	linkAccessory := slack.NewAccessory(button)

	contextSection := slack.NewSectionBlock(contextHeaderText, nil, linkAccessory)

	headerText := slack.NewTextBlockObject("mrkdwn", d.Title, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	bodyText := slack.NewTextBlockObject("mrkdwn", githubmarkdownconvertergo.Slack(d.Body), false, false)
	bodySection := slack.NewSectionBlock(bodyText, nil, nil)

	return slack.NewBlockMessage(
		contextSection,
		dividerSection,
		headerSection,
		bodySection,
	)
}
