package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	url, err := env("INPUT_URL")
	if err != nil {
		log.Fatal(err)
	}

	method, err := env("INPUT_METHOD")
	if err != nil {
		log.Fatal(err)
	}

	body, err := env("INPUT_BODY")
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	if body != "" {
		stringReader := strings.NewReader(body)
		req.Body = io.NopCloser(stringReader)
	}

	client := http.Client{}

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error making request", err)
		return
	}
	defer response.Body.Close()

}

func env(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("env key %s does not exist", key)
	}
	if len(value) == 0 {
		return "", fmt.Errorf("env key %s is empty", key)
	}
	return value, nil
}
