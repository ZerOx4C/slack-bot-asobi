package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Config struct {
	AppToken string `json:"app-token"`
	BotToken string `json:"bot-token"`
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		panic(err)
	}

	// see: https://qiita.com/seratch/items/c7d9aeb60ead5c126c01
	webApi := slack.New(
		config.BotToken,
		slack.OptionAppLevelToken(config.AppToken),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
	socketMode := socketmode.New(
		webApi,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "sock: ", log.Lshortfile|log.LstdFlags)),
	)
	authTest, err := webApi.AuthTest()
	if err != nil {
		panic(err)
	}
	myUserId := authTest.UserID

	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeEventsAPI:
				socketMode.Ack(*envelope.Request)

				eventPayload, _ := envelope.Data.(slackevents.EventsAPIEvent)
				switch eventPayload.Type {
				case slackevents.CallbackEvent:
					switch event := eventPayload.InnerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						if event.User != myUserId {
							println("何もしませんの")
						}
					case *slackevents.AppMentionEvent:
						_, _, err := webApi.PostMessage(
							event.Channel,
							slack.MsgOptionText("なんですの！？", false),
						)
						if err != nil {
							log.Printf("fail: %v", err)
						}
					}
				}
			}
		}
	}()

	socketMode.Run()
}

func loadConfig(filename string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
