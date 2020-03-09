package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/adjust/rmq"
	"github.com/azuretek/crawler/internal/job"
	"github.com/azuretek/crawler/internal/store"
)

// CrawlConsumer is the consumer object for rmq
type CrawlConsumer struct {
	Client *store.Client
}

// Consume performs the crawling work and publishes results to our data store
func (c *CrawlConsumer) Consume(delivery rmq.Delivery) {
	newURL := &job.URL{}
	if err := json.Unmarshal([]byte(delivery.Payload()), &newURL); err != nil {
		delivery.Reject()
		log.Fatal(err)
	}

	// Load current job state
	j, err := job.Get(c.Client, fmt.Sprintf("jobs/%v", newURL.ID))
	if err != nil {
		log.Fatal(err)
	}

	// Get body from request URL
	body, err := GetBody(newURL.URL)
	if err != nil {
		log.Fatal(err)
	}

	// Count number of hostnames in the result
	hostnames, err := CountHostNames(body)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through results and update store counts
	for _, h := range hostnames {
		key := fmt.Sprintf("results/%v/%v", newURL.ID, h.Name)
		if err := c.Client.Increment(key, h.Count); err != nil {
			log.Fatal(err)
		}

		val, err := c.Client.Redis.Get(key).Result()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(key, ": ", val)
	}

	// Normally I'd re-use my connections but this library has made it a private struct
	// If I had more time I would find a different library or make my own implementation
	connection := rmq.OpenConnection("crawler", "tcp", "localhost:6379", 1)
	queue := connection.OpenQueue("crawl_jobs")

	// Create url crawl requests
	urls, err := GetURLs(newURL, body)
	if err != nil {
		log.Fatal(err)
	}

	for _, url := range urls {
		// Only add new crawl jobs if depth is below MaxDepth for this job
		if url.Depth < j.MaxDepth {
			urlBytes, err := json.Marshal(url)
			if err != nil {
				log.Fatal(err)
			}
			queue.PublishBytes(urlBytes)
			log.Printf("Created crawl job for: %v (%v)", url.URL, url.ID)
		}
	}

	delivery.Ack()
}
