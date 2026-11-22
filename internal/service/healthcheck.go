package service

//go:generate mockery --name HealthCheck -- filename healthcheck.go
// HealthCheck defines the interface for the health check service.
type HealthCheck interface {
	GetInstanceID() string
	GetServiceName() string
}

// healthCheckService implements the HealthCheck interface.
type healthCheckService struct {
	serviceName string
	instanceID  string
}

// NewHealthCheck creates a new instance of healthCheckService with the provided service name and instance ID.
func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	return &healthCheckService{
		serviceName: serviceName,
		instanceID:  instanceID,
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
