package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ntduc6403/bookmark-management/internal/repository/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

)

func TestShortenUrl(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name      string
		url       string
		expTime   int
		genCode   func(length int) (string, error)
		setupMock func(ctx context.Context) *mocks.URLStorage
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
			setupMock: func(ctx context.Context) *mocks.URLStorage {
				repo := mocks.NewURLStorage(t)
				repo.On("GetURL", ctx, "AbC123x").Return("", redis.Nil)
				repo.On("StoreURL", ctx, "AbC123x", targetURL, defaultURLExp).Return(nil)
				return repo
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
			setupMock: func(ctx context.Context) *mocks.URLStorage {
				repo := mocks.NewURLStorage(t)
				repo.On("GetURL",ctx, "XyZ987w").Return("", redis.Nil)
				repo.On("StoreURL", ctx, "XyZ987w", targetURL, 3600*time.Second).Return(nil)
				return repo
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
			setupMock: func(ctx context.Context) *mocks.URLStorage {
				repo := mocks.NewURLStorage(t)
				repo.On("GetURL", ctx, "DupCode").Return("https://existing-url.com", nil)
				repo.On("GetURL", ctx, "NewCode").Return("", redis.Nil)
				repo.On("StoreURL", ctx, "NewCode", targetURL, 600*time.Second).Return(nil)
				return repo
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
			setupMock: func(ctx context.Context) *mocks.URLStorage {

				return nil
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
			setupMock: func(ctx context.Context) *mocks.URLStorage {
				repo := mocks.NewURLStorage(t)
				repo.On("GetURL", ctx, "ErrCode").Return("", dbErr)
				return repo
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
			setupMock: func(ctx context.Context) *mocks.URLStorage {
				repo := mocks.NewURLStorage(t)
				repo.On("GetURL", ctx, "SaveErr").Return("", redis.Nil)
				repo.On("StoreURL",ctx, "SaveErr", targetURL, defaultURLExp).Return(dbErr)

				return repo
			},
			wantCode: " ",
			wantErr:  dbErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockRepo := tc.setupMock(ctx)
			svc := &shortenUrl{
				repo:    mockRepo,
				genCode: tc.genCode,
			}

			gotCode, gotErr := svc.ShortenUrl(ctx, tc.url, tc.expTime)

			if tc.wantErr != nil {
				assert.ErrorIs(t, gotErr, tc.wantErr)
			} else {
				assert.NoError(t, gotErr)
			}
			assert.Equal(t, tc.wantCode, gotCode)
		})
	}
}
