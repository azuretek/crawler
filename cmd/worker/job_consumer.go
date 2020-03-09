package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/adjust/rmq"
	"github.com/azuretek/crawler/internal/job"
	"github.com/azuretek/crawler/internal/store"
)

// JobConsumer is the consumer object for rmq
type JobConsumer struct {
	Client *store.Client
}

// Consume looks for new jobs and publishes work to our data store
func (c *JobConsumer) Consume(delivery rmq.Delivery) {
	newJob := &job.Job{}
	if err := json.Unmarshal([]byte(delivery.Payload()), &newJob); err != nil {
		delivery.Reject()
		log.Fatal(err)
	}

	// Take new job and create object in data store to keep track of state
	newJob.State = "In Progress"
	jobBytes, err := json.Marshal(newJob)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Client.Redis.Set(fmt.Sprintf("jobs/%v", newJob.ID), string(jobBytes), 0).Err(); err != nil {
		log.Fatal(err)
	}

	// Normally I'd re-use my connections but this library has made it a private struct
	// If I had more time I would find a different library or make my own implementation
	connection := rmq.OpenConnection("crawler", "tcp", "localhost:6379", 1)
	queue := connection.OpenQueue("crawl_jobs")

	// Take new job and create a url crawl request
	urlBytes, err := json.Marshal(job.URL{
		URL:   newJob.URL,
		ID:    newJob.ID,
		Depth: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
	queue.PublishBytes(urlBytes)
	log.Printf("Created crawl job for: %v (%v)", newJob.URL, newJob.ID)

	delivery.Ack()
}
