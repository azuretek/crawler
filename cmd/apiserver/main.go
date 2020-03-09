package main

import (
	"encoding/json"
	"log"

	"github.com/adjust/rmq"
	"github.com/azuretek/crawler/internal/job"
	uuid "github.com/satori/go.uuid"
)

func main() {
	connection := rmq.OpenConnection("crawler", "tcp", "localhost:6379", 1)
	queue := connection.OpenQueue("new_jobs")

	newJob := &job.Job{
		URL:      "http://example.com",
		ID:       uuid.NewV4(),
		MaxDepth: 2,
	}

	jobBytes, err := json.Marshal(*newJob)
	if err != nil {
		log.Fatal(err)
	}
	queue.PublishBytes(jobBytes)
	log.Printf("New Job Created: %v", newJob.ID)
}
