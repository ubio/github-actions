package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

var (
	vars EnvVars
	err  error
)

// EnvVars passed by GH actions
type EnvVars struct {
	URL    string `envconfig:"INPUT_URL" required:"true"`
	Body   string `envconfig:"INPUT_BODY" required:"true"`
	Secret string `envconfig:"INPUT_SECRET" required:"false"`
}

func main() {

	err = envconfig.Process("", &vars)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	r, err := http.NewRequest("POST", vars.URL, bytes.NewReader([]byte(vars.Body)))
	if err != nil {
		log.Fatal(err)
	}

	if vars.Secret != "" {
		r.Header.Add("X-ORIGIN-SECRET", vars.Secret)
	}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println("::set-output name=status::", string(resp.StatusCode))

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("::set-output name=body::", string(bodyBytes))
}
