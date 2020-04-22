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

func (c cert) slackBlock() *slack.SectionBlock {
	return slack.NewSectionBlock(
		nil,
		[]*slack.TextBlockObject{
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf(":globe_with_meridians: %s\n", c.DomainName), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf(":suspect: %s\n", c.Issuer), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf(":point_right: %s\n", c.IP), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf(":bomb: %d days\n", c.until()), false, false),
		},
		nil,
	)
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
		fmt.Println(cert.DomainName, "expires in", cert.until(), "days")
	}

	fields := make([]slack.Block, 0)
	fields = append(
		fields,
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Our SSL Certificates have been checked and need attention:", false, false),
			nil,
			nil,
		),
	)

	div := slack.NewDividerBlock()
	for i, cert := range certs {
		if i == 0 {
			fields = append(fields, cert.slackBlock())
		} else {
			fields = append(fields, div, cert.slackBlock())
		}
	}

	api := slack.New(os.Getenv("INPUT_SLACK_TOKEN"))
	_, _, err := api.PostMessage("@aw", slack.MsgOptionBlocks(fields...))
	if err != nil {
		log.Fatal(err)
	}

}
