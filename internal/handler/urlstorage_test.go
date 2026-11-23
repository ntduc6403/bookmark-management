package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ntduc6403/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		requestBody      string
		setupMockService func(ctx context.Context) *mocks.ShortenUrl
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Successful",
			requestBody: `{"url":"https://google.com","exp":60}`,
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				mockService := &mocks.ShortenUrl{}
				mockService.On(
					"ShortenUrl",
					ctx,
					"https://google.com",
					60,
				).Return("abc123", nil)

				return mockService
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"code":"abc123","message":"Shorten URL generated successfully!"}`,
		},
		{
			name:        "Invalid request body",
			requestBody: `{"url":""}`,
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				return &mocks.ShortenUrl{}
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid request body"}`,
		},
		{
			name:        "Invalid URL",
			requestBody: `{"url":"invalid-url","exp":60}`,
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				return &mocks.ShortenUrl{}
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid request body"}`,
		},
		{
			name:        "Negative expiration time",
			requestBody: `{"url":"https://google.com","exp":-1}`,
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				return &mocks.ShortenUrl{}
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Expiration time must not be negative"}`,
		},
		{
			name:        "Service error",
			requestBody: `{"url":"https://google.com","exp":60}`,
			setupMockService: func(ctx context.Context) *mocks.ShortenUrl {
				mockService := &mocks.ShortenUrl{}
				mockService.On(
					"ShortenUrl",
					ctx,
					"https://google.com",
					60,
				).Return("", testErr)

				return mockService
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"message":"Failed to shorten URL"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			ctx.Request = httptest.NewRequest(
				http.MethodPost,
				"/shorten",
				strings.NewReader(tc.requestBody),
			)
			ctx.Request.Header.Set("Content-Type", "application/json")

			mockSvc := tc.setupMockService(ctx)
			testHandler := NewShortenURLHandler(mockSvc)

			testHandler.Shorten(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}
