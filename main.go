package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/robfig/cron/v3"
)

func main() {
	// set default log output to stdout
	log.SetOutput(os.Stdout)

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

	// schedule a job to check the health of the sites
	c.AddFunc(config.HealthCheckCronExp, func() {
		for i := 0; i < len(config.Sites); i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				crawl(config, index)
			}(i)
		}
	})

	// schedule a job to check the potential leaked sensitive files of the website
	c.AddFunc(config.SensitiveFilesCronExp, func() {
		for i := 0; i < len(config.Sites); i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				checkSensitiveFiles(config, index)
			}(i)
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
	wg.Wait()
	log.Println("done")
}
