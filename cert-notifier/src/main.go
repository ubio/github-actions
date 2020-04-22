package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nlopes/slack"
	"github.com/prometheus/common/log"
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

	attachments := make([]slack.Attachment, 0)

	attachments = append(attachments, slack.Attachment{
		Color:    "good",
		Fallback: "You successfully posted by Incoming Webhook URL!",
		Text:     "Text in Slack uses the same system of escaping: chat messages, direct messages, file comments, etc. :smile:\nSee <https://api.slack.com/docs/message-formatting#linking_to_channels_and_users>",
	})
	msg := slack.WebhookMessage{
		Attachments: attachments,
	}

	err := slack.PostWebhook(os.Getenv("INPUT_SLACK_URL"), &msg)
	if err != nil {
		fmt.Println(err)
	}
}
