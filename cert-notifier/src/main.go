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

	api := slack.New(os.Getenv("INPUT_SLACK_TOKEN"))
	attachment := slack.Attachment{
		Pretext: "some pretext",
		Text:    "some text",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "a",
				Value: "no",
			},
		},
	}

	channelID, timestamp, err := api.PostMessage("@aw", slack.MsgOptionText("Some text", false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
