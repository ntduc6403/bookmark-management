package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ntduc6403/bookmark-management/internal/api"
	"github.com/ntduc6403/bookmark-management/pkg/redis"
	"github.com/stretchr/testify/assert"
)

// TestShortenLinkEndPoint tests the shorten link endpoint of the API.
func TestShortenLinkEndPoint(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		name             string
		setupTestHTTP    func(engine api.Engine) *httptest.ResponseRecorder
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				body := []byte(`{"url":"https://example.com/long-url","exp":3600}`)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				respRecorder := httptest.NewRecorder()
				engine.ServeHTTP(respRecorder, req)

				return respRecorder
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `"message":"Shorten URL generated successfully!"`,
		},
		{
			name: "invalid body - missing or invalid url",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				body := []byte(`{"url":"invalid-url","exp":3600}`)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				respRecorder := httptest.NewRecorder()
				engine.ServeHTTP(respRecorder, req)

				return respRecorder
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `"message":"Invalid request body"`,
		},
		{
			name: "negative exp time",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				body := []byte(`{"url":"https://example.com/long-url","exp":-1}`)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				respRecorder := httptest.NewRecorder()
				engine.ServeHTTP(respRecorder, req)

				return respRecorder
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `"message":"Expiration time must not be negative"`,
		},
		{
			name: "wrong endpoint",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/shorten", nil)
				respRecorder := httptest.NewRecorder()
				engine.ServeHTTP(respRecorder, req)

				return respRecorder
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: ``,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testAPI := api.NewEngine(
				&api.Config{ServiceName: "bookmark", InstanceID: "123"},
				redis.InitMockRedis(t),
			)

			recorder := tc.setupTestHTTP(testAPI)
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			if tc.expectedResponse != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
			}
		})
	}
}
