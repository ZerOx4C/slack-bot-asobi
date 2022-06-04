package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type config struct {
	GasUrl string `json:"gas-url"`
}

type request struct {
	uuid        string
	unixtime    int64
	channelId   string
	userId      string
	command     string
	text        string
	responseUrl string
}

type slackPayload struct {
	Text         *string `json:"text,omitempty"`
	ResponseType *string `json:"response_type,omitempty"`
}

func str(s string) *string {
	return &s
}

func post(url *url.URL, values url.Values) (string, error) {
	response, err := http.PostForm(url.String(), values)
	if err != nil {
		return "", err
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

func jsonStringify(object any) (string, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
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

	println("running...")
	for {
		err = polling(config)
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
	}
}

func polling(config config) error {
	response, err := http.Get(config.GasUrl)
	if err != nil {
		return err
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	queue := [][]any{}
	err = json.Unmarshal(responseBytes, &queue)
	if err != nil {
		return err
	}

	for _, stringQueueItem := range queue {
		request := request{}
		request.uuid = stringQueueItem[0].(string)
		request.unixtime, _ = stringQueueItem[1].(int64)
		request.channelId = stringQueueItem[2].(string)
		request.userId = stringQueueItem[3].(string)
		request.command = stringQueueItem[4].(string)
		request.text = stringQueueItem[5].(string)
		request.responseUrl = stringQueueItem[6].(string)

		err := process(config, request)
		if err != nil {
			return err
		}
	}

	return nil
}

func process(config config, request request) error {
	slackResponseUrl, err := url.Parse(request.responseUrl)
	if err != nil {
		return err
	}

	slackResponseValues := url.Values{}

	if request.command == "/sushi" {
		slackPayload := slackPayload{}
		slackPayload.Text = str(fmt.Sprintf("hey! %s omachi!", request.text))
		slackPayload.ResponseType = str("in_channel")
		payloadJson, err := jsonStringify(slackPayload)
		if err != nil {
			return err
		}

		slackResponseValues.Set("payload", payloadJson)
	}

	response, err := post(slackResponseUrl, slackResponseValues)
	if err != nil {
		return err
	}

	println(response)

	gasDoneUrl, err := url.Parse(config.GasUrl)
	if err != nil {
		return err
	}

	gasDoneValues := url.Values{}
	gasDoneValues.Set("uuid", request.uuid)

	response, err = post(gasDoneUrl, gasDoneValues)
	if err != nil {
		return err
	}

	return nil
}
