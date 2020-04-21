package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

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

func main() {
	input := os.Getenv("INPUT_CERTS")
	warnUnderDays, err := strconv.ParseFloat(
		os.Getenv("INPUT_WARN_UNDER_DAYS"),
		64,
	)
	if err != nil {
		log.Fatal(err)
	}

	certs := make([]cert, 0)
	if err := json.Unmarshal([]byte(input), &certs); err != nil {
		log.Fatal(err)
	}

	l := "2006-01-02 15:04:05 -0700 MST"
	now := time.Now()
	for _, cert := range certs {
		expires, err := time.Parse(l, cert.NotAfter)
		if err != nil {
			log.Fatal(err)
		}

		days := expires.Sub(now).Hours() / 24
		fmt.Printf("Checked %s. Expires in %f days", cert.DomainName, math.Round(days))
		if days < warnUnderDays {
			fmt.Println(cert.DomainName, "expiring!!! In", days, "days")
		}
	}
}
