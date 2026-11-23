package repository

import (
	"context"
	"testing"

	redisPkg "github.com/ntduc6403/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckRepo_Ping(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name      string
		setupMock func() *redis.Client
		expectErr bool
	}{
		{
			name: "normal case - redis is alive",
			setupMock: func() *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectErr: false,
		},
		{
			name: "error case - redis client closed",
			setupMock: func() *redis.Client {
				client := redisPkg.InitMockRedis(t)
				_ = client.Close()
				return client
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			redisMock := tc.setupMock()
			testRepo := NewHealthCheck(redisMock)

			err := testRepo.Ping(ctx)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}