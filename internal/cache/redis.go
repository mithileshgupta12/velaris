package cache

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

type Stores struct {
	SessionStore SessionStore
}

func NewRedisClient() (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	statusCmd := client.Ping(context.Background())
	err := statusCmd.Err()
	if err != nil {
		return nil, err
	}

	return &RedisClient{client}, nil
}

func (rc *RedisClient) InitStores() *Stores {
	sessionStore := NewSessionStore(rc.client)

	return &Stores{sessionStore}
}

func (rc *RedisClient) Close() {
	err := rc.client.Close()
	if err != nil {
		slog.Error("failed to close redis connection", "err", err)
	}
}
