package blacklist

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

var ctx = context.Background()

type Manager struct {
	redisClient *redis.Client
}

func NewBlacklistManager(redisClient *redis.Client) *Manager {
	return &Manager{
		redisClient: redisClient,
	}
}

func (m *Manager) BlacklistPlace(id int) error {
	cmd := m.redisClient.Append(ctx, strconv.FormatInt(int64(id), 10), "")
	_, err := cmd.Result()

	return err
}

func (m *Manager) IsBlacklisted(id int) (bool, error) {
	cmd := m.redisClient.Exists(ctx, strconv.FormatInt(int64(id), 10))
	return cmd.Result()

	return m.redisClient.Exists(ctx, strconv.FormatInt(int64(id), 10)).Result()
}
