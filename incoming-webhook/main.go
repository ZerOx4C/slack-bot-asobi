package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type config struct {
	Url string `json:"url"`
}

type payload struct {
	Text      *string `json:"text,omitempty"`
	UserName  *string `json:"username,omitempty"`
	IconUrl   *string `json:"icon_url,omitempty"`
	IconEmoji *string `json:"icon_emoji,omitempty"`
	Channel   *string `json:"channel,omitempty"`
}

func str(s string) *string {
	return &s
}

func main() {
	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	config := config{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	payload := payload{}
	payload.Text = str("hello.")
	payloadJsonString, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	response, err := http.Post(config.Url, "application/json", bytes.NewBuffer(payloadJsonString))
	if err != nil {
		panic(err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	println(string(responseBytes))
}
