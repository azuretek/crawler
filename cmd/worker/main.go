package main

import (
	"fmt"
	"time"

	"github.com/adjust/rmq"
	"github.com/azuretek/crawler/internal/store"
	uuid "github.com/satori/go.uuid"
)

func main() {
	connection := rmq.OpenConnection("crawler", "tcp", "localhost:6379", 1)
	client := store.NewClient()

	// Process new job requests
	newJobQueue := connection.OpenQueue("new_jobs")
	newJobQueue.StartConsuming(10, time.Second)
	jobConsumer := &JobConsumer{Client: client}
	newJobQueue.AddConsumer(fmt.Sprint(uuid.NewV4()), jobConsumer)

	// Process URL crawl requests
	crawlQueue := connection.OpenQueue("crawl_jobs")
	crawlQueue.StartConsuming(10, time.Second)
	crawlConsumer := &CrawlConsumer{Client: client}
	crawlQueue.AddConsumer(fmt.Sprint(uuid.NewV4()), crawlConsumer)

	select {}
}
