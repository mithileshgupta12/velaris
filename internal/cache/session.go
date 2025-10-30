package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionStore interface {
	Set(ctx context.Context, key, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type sessionStore struct {
	client *redis.Client
}

func NewSessionStore(client *redis.Client) SessionStore {
	return &sessionStore{client}
}

func (ss *sessionStore) Set(ctx context.Context, key, value any, expiration time.Duration) error {
	return ss.client.Set(ctx, fmt.Sprintf("session:%s", key), value, expiration).Err()
}

func (ss *sessionStore) Get(ctx context.Context, key string) (string, error) {
	result, err := ss.client.Get(ctx, fmt.Sprintf("session:%s", key)).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (ss *sessionStore) Del(ctx context.Context, key string) error {
	return ss.client.Del(ctx, fmt.Sprintf("session:%s", key)).Err()
}
