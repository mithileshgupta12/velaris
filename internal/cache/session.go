package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionStore interface {
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type sessionStore struct {
	client *redis.Client
}

func NewSessionStore(client *redis.Client) SessionStore {
	return &sessionStore{client}
}

func (ss *sessionStore) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	err := ss.client.Set(ctx, fmt.Sprintf("session:%s", key), value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ss *sessionStore) Get(ctx context.Context, key string) (string, error) {
	result, err := ss.client.Get(ctx, fmt.Sprintf("session:%s", key)).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
