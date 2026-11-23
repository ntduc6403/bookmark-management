package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name URLStorage -- filename urlstorage.go
type URLStorage interface {
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
	GetURL(ctx context.Context, code string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{
		c: c,
	}
}
//	StoreURL stores the given code and URL in Redis with a specified expiration time. It returns any error encountered during the storage process.
func (s *urlStorage) StoreURL(ctx context.Context, code, url string , exp time.Duration) error {
	return s.c.Set(ctx, code, url,  exp).Err()

}
// GetURL retrieves the original URL associated with the given code from Redis. It returns the URL and any error encountered during the retrieval process.
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}
