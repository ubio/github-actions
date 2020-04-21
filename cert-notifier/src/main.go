package main

import (
	"encoding/json"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/common/log"
)

type cert struct {
	Domain     string   `json:"domainName"`
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
	spew.Dump(certs)
}
