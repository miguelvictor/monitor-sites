package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func crawl(config Config, index int) {
	// get status code of site url
	response, err := http.Get(config.Sites[index].Site)
	if err != nil {
		log.Println(err)
		sendMail(config, index, "Site Unreachable")
		return
	}
	defer response.Body.Close()

	// send an email if status code is not 200
	if response.StatusCode != 200 {
		log.Printf("[%s]: %s\n", config.Sites[index].Site, response.Status)
		sendMail(config, index, response.Status)
	}
}

func checkSensitiveFiles(config Config, index int) {
	// get http client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// check if some sensitive files are exposed
	paths := [3]string{".env", "php.ini", "wp-config.php"}
	for i := 0; i < len(paths); i++ {
		// visit site with added path
		response, err := client.Get(config.Sites[index].Site + "/" + paths[i])
		if err != nil {
			log.Println(err)
			continue
		}
		defer response.Body.Close()

		// read response body and close it afterwards
		body, err := io.ReadAll(response.Body)

		// check response for errors
		if err != nil {
			log.Printf("Failed to read response body: %s\n", err)
			continue
		}
		if len(strings.TrimSpace(string(body))) != 0 {
			sendSlackNotification(config, index, paths[i], string(body))
		}
	}
}
