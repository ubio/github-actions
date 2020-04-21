package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("Checked %s. Expires in %d days", cert.DomainName, expires)
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
	for _, cert := range certs {
		fmt.Println(cert.DomainName, "expiring!!! In", cert.until(), "days")
	}
}
