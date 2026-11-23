package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ntduc6403/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

var testErr = errors.New("test error")

// TestHealthCheckHandler tests the HealthCheck handler for various scenarios.
func TestHealthCheckHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.HealthCheck
		expectedStatus   int
		expectedResponse string

		shouldGenerateID bool
	}{

		{
			name: "Successful",
			setupRequest: func(ctx *gin.Context) {

				ctx.Request = httptest.NewRequest("GET", "/health-check", nil)
			},

			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				mockService := mocks.NewHealthCheck(t)
				mockService.On("CheckRedis", ctx).Return(nil)
				mockService.On("GetServiceName").Return("bookmark_service")
				mockService.On("GetInstanceID").Return("text-id")
				return mockService
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"OK","service_name":"bookmark_service","instance_id":"text-id"}`,
		},
		{
			name: "Redis check failed",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				mockService := mocks.NewHealthCheck(t)
				mockService.On("CheckRedis", ctx).Return(testErr)
				return mockService
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message":"Internal Server Error"}`,
		},
	}
	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupRequest(ctx)

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewHealthCheck(mockSvc)
			testHandler.HealthCheck(ctx)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
