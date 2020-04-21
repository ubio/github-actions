package main

import (
	"encoding/json"
	"fmt"
	"os"
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
		fmt.Println(cert.DomainName, "expires in", days, "days")
	}
}
