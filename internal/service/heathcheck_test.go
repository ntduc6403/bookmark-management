package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultServiceName string = "bookmark_service"

func TestGetInstanceID(t *testing.T) {
	tests := []struct {
		name string

		expectedId     string
		expectedError  error
		shouldGenerate bool
	}{
		{
			name:          "Existing instance ID",
			expectedId:    "test-id",
			expectedError: nil,
		},
		{
			name:           "Empty instance ID",
			expectedId:     "",
			expectedError:  nil,
			shouldGenerate: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := NewHealthCheck(defaultServiceName, tc.expectedId)
			instanceId, err := svc.GetInstanceID()
			assert.ErrorIs(t, err, tc.expectedError)
			if tc.shouldGenerate {
				assert.NotEmpty(t, instanceId)
				assert.NotEqual(t, tc.expectedId, instanceId)
			} else {
				assert.Equal(t, tc.expectedId, instanceId)
			}

		})
	}
}
