package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	ID         int         `json:"id,omitempty"`
	Title      string      `json:"title,omitempty"`
	Body       string      `json:"body,omitempty"`
	CreatedAt  string      `json:"created_at,omitempty"`
	HTMLURL    string      `json:"html_url,omitempty"`
	Category   *Category   `json:"category,omitempty"`
	User       *User       `json:"user,omitempty"`
	Repository *Repository `json:"repository,omitempty"`
}

type Repository struct {
	Name string `json:"name"`
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
		log.Fatal("error with env:", err)
	}

	file, err := os.Open(vars.EventPath)
	if err != nil {
		log.Fatal("error with event file:", err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("error with content:", err)
	}

	event := DiscussionEvent{}
	err = json.Unmarshal(content, &event)
	if err != nil {
		log.Fatal("error decoding discussion:", err)
	}

	api := slack.New(vars.Token)
	channelID, timestamp, err := api.PostMessage(
		vars.Channel,
		buildSlackBlock(event.Discussion),
	)
	if err != nil {
		log.Fatal("message failed to send:", err)
	}

	log.Printf("message successfully sent to channel %s at %s", channelID, timestamp)
}

func buildSlackBlock(d *Discussion) slack.MsgOption {

	dividerSection := slack.NewDividerBlock()

	// @TODO: get the squad name from the event
	squadName := getSquadNameFromRepoName(d.Repository.Name)

	viewDiscussionText := "View Discussion"

	headerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("%s - *%s* %s ", squadName, d.Category.Name, d.Category.Emoji), false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	titleText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s*", d.Title), false, false)
	titleSection := slack.NewSectionBlock(titleText, nil, nil)

	bodyText := slack.NewTextBlockObject("mrkdwn", githubmarkdownconvertergo.Slack(d.Body), false, false)
	bodySection := slack.NewSectionBlock(bodyText, nil, nil)

	contextTextAuthor := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Author: <%s|%s>", d.User.HTMLURL, d.User.Login), false, false)
	contextTextLink := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<%s|%s>", d.HTMLURL, viewDiscussionText), false, false)
	contextSection := slack.NewContextBlock("", []slack.MixedElement{contextTextAuthor, contextTextLink}...)

	return slack.MsgOptionBlocks(
		headerSection,
		titleSection,
		bodySection,
		dividerSection,
		contextSection,
	)
}

func getSquadNameFromRepoName(repoName string) string {
	p := strings.Split(repoName, "-")
	squad := p[2:]
	squadParts := make([]string, len(squad))
	for _, p := range squad {
		squadParts = append(squadParts, strings.Title(p))
	}
	squadName := strings.Join(squadParts, " ")

	return squadName
}
