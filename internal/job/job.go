package job

import (
	"encoding/json"

	"github.com/azuretek/crawler/internal/store"
	uuid "github.com/satori/go.uuid"
)

// Job contains the details for a job
type Job struct {
	URL      string    `json:"url"`
	ID       uuid.UUID `json:"uuid"`
	MaxDepth int       `json:"maxdepth"`
	State    string    `json:"state"`
}

// Get returns the job details from the store
func Get(store *store.Client, key string) (*Job, error) {
	result, err := store.Redis.Get(key).Result()
	if err != nil {
		return nil, err
	}

	job := &Job{}
	if err := json.Unmarshal([]byte(result), &job); err != nil {
		return nil, err
	}

	return job, nil
}

// URL contains information about urls
type URL struct {
	URL   string    `json:"url"`
	ID    uuid.UUID `json:"uuid"`
	Depth int       `json:"depth"`
}
