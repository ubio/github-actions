package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/common/log"
	"github.com/slack-go/slack"
)

type cert struct {
	DomainName string   `json:"domainName"`
	IP         string   `json:"ip"`
	Issuer     string   `json:"issuer"`
	CommonName string   `json:"commonName"`
	NotBefore  string   `json:"notBefore"`
	NotAfter   string   `json:"notAfter"`
	Error      string   `json:"error"`
	Sans       []string `json:"sans"`
}

func (c cert) until() int64 {

	l := "2006-01-02 15:04:05 -0700 MST"
	now := time.Now()

	expires, err := time.Parse(l, c.NotAfter)
	if err != nil {
		log.Fatal(err)
	}

	return int64(expires.Sub(now).Hours() / 24)
}

func main() {
	input := os.Getenv("INPUT_CERTS")
	warnUnderDays, err := strconv.ParseInt(os.Getenv("INPUT_WARN_UNDER_DAYS"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	certs := make([]cert, 0)
	if err := json.Unmarshal([]byte(input), &certs); err != nil {
		log.Fatal(err)
	}

	expiring := make([]cert, 0)
	for _, cert := range certs {
		expires := cert.until()
		fmt.Printf("Checked %s - expires in %d days (%s)\n", cert.DomainName, expires, cert.NotAfter)
		if expires < warnUnderDays {
			expiring = append(expiring, cert)
		}
	}

	if len(expiring) == 0 {
		return
	}

	warn(expiring)
}

func warn(certs []cert) {
	fmt.Println("-----------------")
	fmt.Println("!!! WARNINGS: !!!")
	fmt.Println("")

	for _, cert := range certs {
		fmt.Println(cert.DomainName, "expiring in", cert.until(), "days")
	}
	headerText := slack.NewTextBlockObject("mrkdwn", "You have a new request:\n*<fakeLink.toEmployeeProfile.com|Fred Enriquez - New device request>*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Fields
	typeField := slack.NewTextBlockObject("mrkdwn", "*Type:*\nComputer (laptop)", false, false)
	whenField := slack.NewTextBlockObject("mrkdwn", "*When:*\nSubmitted Aut 10", false, false)
	lastUpdateField := slack.NewTextBlockObject("mrkdwn", "*Last Update:*\nMar 10, 2015 (3 years, 5 months)", false, false)
	reasonField := slack.NewTextBlockObject("mrkdwn", "*Reason:*\nAll vowel keys aren't working.", false, false)
	specsField := slack.NewTextBlockObject("mrkdwn", "*Specs:*\n\"Cheetah Pro 15\" - Fast, really fast\"", false, false)

	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, typeField)
	fieldSlice = append(fieldSlice, whenField)
	fieldSlice = append(fieldSlice, lastUpdateField)
	fieldSlice = append(fieldSlice, reasonField)
	fieldSlice = append(fieldSlice, specsField)

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

	// Approve and Deny Buttons
	approveBtnTxt := slack.NewTextBlockObject("plain_text", "Approve", false, false)
	approveBtn := slack.NewButtonBlockElement("", "click_me_123", approveBtnTxt)

	denyBtnTxt := slack.NewTextBlockObject("plain_text", "Deny", false, false)
	denyBtn := slack.NewButtonBlockElement("", "click_me_123", denyBtnTxt)

	actionBlock := slack.NewActionBlock("", approveBtn, denyBtn)

	// Build Message with blocks created above

	msg := slack.NewBlockMessage(
		headerSection,
		fieldsSection,
		actionBlock,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(string(b))

	api := slack.New(os.Getenv("INPUT_SLACK_TOKEN"))
	_, _, err = api.PostMessage("@aw", slack.MsgOptionText(string(b), true))

	// attachment := slack.Attachment{
	// 	Pretext: "some pretext",
	// 	Text:    "some text",
	// 	Fields: []slack.AttachmentField{
	// 		slack.AttachmentField{
	// 			Title: "a",
	// 			Value: "no",
	// 		},
	// 	},
	// }

	// channelID, timestamp, err := api.PostMessage("@aw", slack.MsgOptionText("Some text", false), slack.MsgOptionAttachments(attachment))
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// 	return
	// }
	// fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
