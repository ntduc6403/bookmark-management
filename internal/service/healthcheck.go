package service

import (
	"context"

	"github.com/ntduc6403/bookmark-management/internal/repository"
)

// HealthCheck defines the interface for the health check service.
//
//go:generate mockery --name HealthCheck -- filename healthcheck.go
type HealthCheck interface {
	GetInstanceID() string
	GetServiceName() string
	CheckRedis(ctx context.Context) error
}

// healthCheckService implements the HealthCheck interface.
type healthCheckService struct {
	serviceName string
	instanceID  string
	repo        repository.HealthCheck
}

// NewHealthCheck creates a new instance of healthCheckService with the provided service name and instance ID.
func NewHealthCheck(
	serviceName string,
	instanceID string,
	repo repository.HealthCheck,
) HealthCheck {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
		repo:        repo,
	}
}

// GetInstanceID returns the instance ID of the service.

func (h *healthCheckService) GetInstanceID() string {
	return h.instanceID
}

// GetServiceName returns the name of the service.
func (h *healthCheckService) GetServiceName() string {
	return h.serviceName
}
func (h *healthCheckService) CheckRedis(ctx context.Context) error {
	return h.repo.Ping(ctx)
}
