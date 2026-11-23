package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ntduc6403/bookmark-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			svc := NewHealthCheck(defaultServiceName, tt.expectedInstanceID, nil)

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
			svc := NewHealthCheck(tt.expectedServiceName, " ", nil)

			serviceName := svc.GetServiceName()

			assert.Equal(t, tt.expectedServiceName, serviceName)
		})
	}
}

func TestCheckRedis(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		pingErr   error
		expectErr bool
	}{
		{
			name:      "Ping Redis success",
			pingErr:   nil,
			expectErr: false,
		},
		{
			name:      "Ping Redis error",
			pingErr:   errors.New("redis connection failed"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := mocks.NewHealthCheck(t)
			mockRepo.On("Ping", mock.Anything).Return(tt.pingErr).Once()

			svc := NewHealthCheck(defaultServiceName, "inst-1", mockRepo)
			err := svc.CheckRedis(context.Background())

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
