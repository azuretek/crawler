package store

import (
	"github.com/go-redis/redis/v7"
)

// Client contains store client information
type Client struct {
	Redis *redis.Client
}

// NewClient returns a new store client
func NewClient() *Client {
	return &Client{
		Redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

// Increment bumps the count for the specified key in a tranasctionally safe way
func (c *Client) Increment(key string, count int64) error {
	pipe := c.Redis.TxPipeline()
	if err := pipe.IncrBy(key, count).Err(); err != nil {
		return err
	}

	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}
