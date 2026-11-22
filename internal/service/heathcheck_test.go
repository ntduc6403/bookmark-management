package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultServiceName string = "bookmark_service"

func TestGetInstanceID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		expectedInstanceID string
		want               string
	}{
		{
			name:               "Valid InstanceID",
			expectedInstanceID: "123",
		},
		{
			name:               "Empty InstanceID",
			expectedInstanceID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewHealthCheck(defaultServiceName, tt.expectedInstanceID)

			instanceID := svc.GetInstanceID()

			assert.Equal(t, tt.expectedInstanceID, instanceID)
		})
	}
}
func TestGetServiceName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                string
		expectedServiceName string
	}{
		{
			name:                "Valid service name",
			expectedServiceName: "health-check",
		},
		{
			name:                "Empty service name",
			expectedServiceName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewHealthCheck(tt.expectedServiceName, " ")

			serviceName := svc.GetServiceName()

			assert.Equal(t, tt.expectedServiceName, serviceName)
		})
	}
}
