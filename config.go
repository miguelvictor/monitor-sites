package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	SmtpUser              string   `json:"smtpUser"`
	SmtpPassword          string   `json:"smtpPassword"`
	Emails                []string `json:"emails"`
	HealthCheckCronExp    string   `json:"healthCheckCron"`
	SensitiveFilesCronExp string   `json:"sensitiveFilesCron"`
	SlackWebhook          string   `json:"slackWebhook"`
	Sites                 []struct {
		Site   string   `json:"site"`
		Emails []string `json:"emails"`
	} `json:"sites"`
}

func getSitesConfig() (config Config, err error) {
	// read json file called config.json
	data, err := os.ReadFile("config.json")
	if err != nil {
		return
	}

	// parse config json file
	err = json.Unmarshal(data, &config)
	if err != nil {
		return
	}

	return config, nil
}
