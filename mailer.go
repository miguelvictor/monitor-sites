package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
)

func sendMail(config Config, index int, response *http.Response) {
	parsedUrl, err := url.Parse(config.Sites[index].Site)
	if err != nil {
		log.Printf("Invalid site url: %s\n", config.Sites[index].Site)
		log.Println(err)
		return
	}

	to := mergeToEmails(config, index)
	auth := smtp.PlainAuth("", config.SmtpUser, config.SmtpPassword, "smtp.gmail.com")
	msg := fmt.Sprintf("Site: %s<br>Status: %s", config.Sites[index].Site, response.Status)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	payload := fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: %s\n%s\n\n%s",
		fmt.Sprintf("DF DevOps <%s>", config.SmtpUser),
		strings.Join(to, ","),
		fmt.Sprintf("%s is down", parsedUrl.Host),
		mime,
		msg,
	)
	smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		config.SmtpUser,
		to,
		[]byte(payload),
	)
}

func mergeToEmails(config Config, index int) []string {
	// no emails in site config
	if len(config.Sites[index].Emails) == 0 {
		return config.Emails
	}

	// merge site emails with global emails
	seen := make(map[string]bool)
	var result []string

	// add elements from config.Emails to result
	for _, email := range config.Emails {
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}

	// add elements from config.Sites[index].Emails to result
	for _, email := range config.Sites[index].Emails {
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}

	return result
}
