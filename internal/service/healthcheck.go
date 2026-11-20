package service

import "github.com/google/uuid"
//go:generate mockery --name HealthCheck -- filename healthcheck.go 
// HealthCheck defines the interface for the health check service.
type HealthCheck interface {
	GetInstanceID() (string, error)
	GetServiceName() string
}
// healthCheckService implements the HealthCheck interface.
type healthCheckService struct {
	serviceName string
	instanceID string
}
// NewHealthCheck creates a new instance of healthCheckService with the provided service name and instance ID.
func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}
// GetInstanceID returns the instance ID of the service. If the instance ID is empty, it generates a new UUID and returns it.
func (h *healthCheckService) GetInstanceID() (string, error) {

	if h.instanceID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}

		h.instanceID = id.String()
	}
	return h.instanceID, nil
}
// GetServiceName returns the name of the service.
func (h *healthCheckService) GetServiceName() string {
	return h.serviceName
}
