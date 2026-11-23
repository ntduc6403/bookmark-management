package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name HealthCheck -- filename healthcheck.go

type HealthCheck interface {
	Ping(ctx context.Context) error
}

type healthCheck struct {
	client *redis.Client
}

func NewHealthCheck(client *redis.Client) HealthCheck {
	return &healthCheck{
		client: client,
	}
}

// Hàm này gọi Redis PING
func (r *healthCheck) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
