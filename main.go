package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/robfig/cron/v3"
)

func crawl(mutex *sync.Mutex, wg *sync.WaitGroup, config Config, index int) {
	// add job to the wait group
	mutex.Lock()
	wg.Add(1)
	mutex.Unlock()

	// get status code of site url
	log.Printf("[%s] checking has started\n", config.Sites[index].Site)
	response, err := http.Get(config.Sites[index].Site)
	if err != nil {
		log.Printf("Failed to fetch site url: %s\n", config.Sites[index].Site)
		log.Println(err)
	}

	// send an email if status code is 200
	log.Printf("[%s]: %s\n", config.Sites[index].Site, response.Status)
	if response.StatusCode != 200 {
		sendMail(config, index, response)
	}

	// remove job from the wait group
	mutex.Lock()
	wg.Done()
	mutex.Unlock()
}

func main() {
	// parse config json file
	config, err := getSitesConfig()
	if err != nil {
		log.Fatal("Failed to parse config.json: ", err)
	}

	// create a new cron job scheduler
	c := cron.New()

	// create a wait group to wait for the running jobs to finish
	// we also create a mutex to lock the wait group
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// schedule a job to run every minute
	c.AddFunc(config.CronExp, func() {
		for i := 0; i < len(config.Sites); i++ {
			go crawl(&mutex, &wg, config, i)
		}
	})

	// start the cron job scheduler
	c.Start()
	log.Println("cron job scheduler started")

	// wait for a termination signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// gracefully stop the cron job scheduler
	c.Stop()

	// wait for the running jobs to stop
	log.Println("waiting for running jobs to finish...")
	mutex.Lock()
	wg.Wait()
	mutex.Unlock()
	log.Println("done")
}
