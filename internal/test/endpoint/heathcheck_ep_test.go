package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ntduc6403/bookmark-management/internal/api"
	"github.com/stretchr/testify/assert"
)

// TestHealthCheckEndPoint tests the health check endpoint of the API.
func TestHealthCheckEndPoint(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name             string
		setupTestHTTP    func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("GET", "/health-check", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)

				return respRecorder
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `"message":"OK","service_name":"","instance_id":`,
		},
		{
			name: "wrong endpoint",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest("POST", "/health-check", nil)
				respRecorder := httptest.NewRecorder()
				api.ServeHTTP(respRecorder, req)

				return respRecorder
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: ``,
		},
	}
	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			testAPI := api.NewEngine(&api.Config{})

			recorder := tc.setupTestHTTP(testAPI)
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
		})
	}
}
