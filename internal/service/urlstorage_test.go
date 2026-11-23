package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ntduc6403/bookmark-management/internal/repository/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenUrl(t *testing.T) {
	type testCase struct {
		name      string
		url       string
		expTime   int
		genCode   func(length int) (string, error)
		setupMock func(repo *mocks.URLStorage)
		wantCode  string
		wantErr   error
	}

	targetURL := "https://example.com/very-long-bookmark-url"
	dbErr := errors.New("redis connection failed")
	genErr := errors.New("rand failed")

	tests := []testCase{
		{
			name:    "success with default expiration when expTime <= 0",
			url:     targetURL,
			expTime: 0,
			genCode: func(length int) (string, error) {
				assert.Equal(t, urlCodeLength, length)
				return "AbC123x", nil
			},
			setupMock: func(repo *mocks.URLStorage) {
				repo.On("GetURL", mock.Anything, "AbC123x").
					Return("", redis.Nil).
					Once()
				repo.On("StoreURL", mock.Anything, "AbC123x", targetURL, defaultURLExp).
					Return(nil).
					Once()
			},
			wantCode: "AbC123x",
			wantErr:  nil,
		},
		{
			name:    "success with custom expiration when expTime > 0",
			url:     targetURL,
			expTime: 3600,
			genCode: func(length int) (string, error) {
				assert.Equal(t, urlCodeLength, length)
				return "XyZ987w", nil
			},
			setupMock: func(repo *mocks.URLStorage) {
				repo.On("GetURL", mock.Anything, "XyZ987w").
					Return("", redis.Nil).
					Once()
				repo.On("StoreURL", mock.Anything, "XyZ987w", targetURL, 3600*time.Second).
					Return(nil).
					Once()
			},
			wantCode: "XyZ987w",
			wantErr:  nil,
		},
		{
			name:    "collision on first code then retry returns second code",
			url:     targetURL,
			expTime: 600,
			genCode: func() func(length int) (string, error) {
				codes := []string{"DupCode", "NewCode"}
				idx := 0
				return func(length int) (string, error) {
					code := codes[idx]
					idx++
					return code, nil
				}
			}(),
			setupMock: func(repo *mocks.URLStorage) {
				repo.On("GetURL", mock.Anything, "DupCode").
					Return("https://existing-url.com", nil).
					Once()
				repo.On("GetURL", mock.Anything, "NewCode").
					Return("", redis.Nil).
					Once()
				repo.On("StoreURL", mock.Anything, "NewCode", targetURL, 600*time.Second).
					Return(nil).
					Once()
			},
			wantCode: "NewCode",
			wantErr:  nil,
		},
		{
			name:    "error when genCode fails",
			url:     targetURL,
			expTime: 0,
			genCode: func(length int) (string, error) {
				return "", genErr
			},
			setupMock: func(repo *mocks.URLStorage) {
			
			},
			wantCode: " ",
			wantErr:  genErr,
		},
		{
			name:    "error when GetURL returns unexpected error",
			url:     targetURL,
			expTime: 0,
			genCode: func(length int) (string, error) {
				return "ErrCode", nil
			},
			setupMock: func(repo *mocks.URLStorage) {
				repo.On("GetURL", mock.Anything, "ErrCode").
					Return("", dbErr).
					Once()
			},
			wantCode: "",
			wantErr:  dbErr,
		},
		{
			name:    "error when StoreURL fails",
			url:     targetURL,
			expTime: 0,
			genCode: func(length int) (string, error) {
				return "SaveErr", nil
			},
			setupMock: func(repo *mocks.URLStorage) {
				repo.On("GetURL", mock.Anything, "SaveErr").
					Return("", redis.Nil).
					Once()
				repo.On("StoreURL", mock.Anything, "SaveErr", targetURL, defaultURLExp).
					Return(dbErr).
					Once()
			},
			wantCode: " ",
			wantErr:  dbErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := mocks.NewURLStorage(t)
			if tc.setupMock != nil {
				tc.setupMock(mockRepo)
			}

			svc := &shortenUrl{
				repo:    mockRepo,
				genCode: tc.genCode,
			}

			gotCode, gotErr := svc.ShortenUrl(context.Background(), tc.url, tc.expTime)

			if tc.wantErr != nil {
				assert.ErrorIs(t, gotErr, tc.wantErr)
			} else {
				assert.NoError(t, gotErr)
			}
			assert.Equal(t, tc.wantCode, gotCode)
		})
	}
}