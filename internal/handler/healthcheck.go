package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntduc6403/bookmark-management/internal/service"
)

// HealthCheckResponse represents the response structure for the health check endpoint.
type HealthCheckResponse struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

// HealthCheck defines the interface for the health check handler.
type HealthCheck interface {
	HealthCheck(c *gin.Context)
}

// healthCheckHandler implements the HealthCheck interface.
type healthCheckHandler struct {
	HealthCheckService service.HealthCheck
}

// NewHealthCheck creates a new instance of healthCheckHandler with the provided HealthCheck service.
func NewHealthCheck(healthCheckSvc service.HealthCheck) HealthCheck {
	return &healthCheckHandler{
		HealthCheckService: healthCheckSvc,
	}
}

// HealthCheck handles the health check endpoint, returning the service name and instance ID.
// @Summary HealthCheck endpoint
// @Description Returns the health status of the service, including service name and instance ID.
// @Tags HealthCheck
// @Success 200 {object} HealthCheckResponse
// @Failure 500 {object} map[string]string
// @Router /health-check [get]
func (h *healthCheckHandler) HealthCheck(c *gin.Context) {

	err := h.HealthCheckService.CheckRedis(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	serviceName := h.HealthCheckService.GetServiceName()
	instanceID := h.HealthCheckService.GetInstanceID()

	c.JSON(http.StatusOK, HealthCheckResponse{
		Message:     "OK",
		ServiceName: serviceName,
		InstanceID:  instanceID,
	})

}
