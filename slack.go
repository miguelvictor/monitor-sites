package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SlackData struct {
	Text   string           `json:"text"`
	Blocks []SlackDataBlock `json:"blocks"`
}

type SlackDataBlock struct {
	Type string             `json:"type"`
	Text SlackDataBlockText `json:"text"`
}

type SlackDataBlockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func sendSlackNotification(config Config, index int, path string, content string) {
	// noop if slack webhook is not set
	if config.SlackWebhook == "" {
		return
	}

	// parse url
	parsedUrl, err := url.Parse(config.Sites[index].Site)
	if err != nil {
		log.Printf("Invalid site url: %s\n", config.Sites[index].Site)
		log.Println(err)
		return
	}

	// prepare slack payload data
	title := fmt.Sprintf("Exposed file from %s", parsedUrl.Host)
	exposedFileUrl := parsedUrl.JoinPath(path)
	data := SlackData{
		Text: title,
		Blocks: []SlackDataBlock{
			{
				Type: "section",
				Text: SlackDataBlockText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*%s*\nURL: <%s|%s>\n```\n%s\n```", title, exposedFileUrl, exposedFileUrl, strings.TrimSpace(content)),
				},
			},
		},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal slack notification json: ", err)
		return
	}

	// send slack notification
	response, err := http.Post(config.SlackWebhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Failed to send slack notification: ", err)
		return
	}
	defer response.Body.Close()
}
