package repository

import (
	"context"
	"testing"

	redisPkg "github.com/ntduc6403/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestURLStorage_StoreURL(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		setupMock   func() *redis.Client
		expectedErr error
		verifyFunc  func(ctx context.Context, c *redis.Client)
	}{
		{
			name: "normal case",
			setupMock: func() *redis.Client {

				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			expectedErr: nil,
			verifyFunc: func(ctx context.Context, c *redis.Client) {

				url, err := c.Get(ctx, "1234567").Result()
				assert.Nil(t, err)
				assert.Equal(t, url, "https://example.com")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			redisMock := tc.setupMock()
			testRepo := NewURLStorage(redisMock)

			err := testRepo.StoreURL(context.Background(), "1234567", "https://example.com" ,12312)
			assert.Equal(t, tc.expectedErr, err)
			if err == nil {
				tc.verifyFunc(ctx , redisMock)
			}
		})

	}
}
