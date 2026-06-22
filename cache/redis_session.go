package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisSession struct {
	client *redis.Client
}

// NewRedisSession wraps a go-redis client; client stays inside cache package.
func NewRedisSession(c *redis.Client) RedisSession {
	if c == nil {
		return nil
	}
	return &redisSession{client: c}
}

func (s *redisSession) Get(ctx context.Context, key string) (string, error) {
	if s == nil || s.client == nil {
		return "", nil
	}
	val, err := s.client.Get(ctx, key).Result()
	return val, NormalizeErr(err)
}

func (s *redisSession) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if s == nil || s.client == nil {
		return nil
	}
	err := s.client.Set(ctx, key, value, ttl).Err()
	return NormalizeErr(err)
}

func (s *redisSession) Del(ctx context.Context, keys ...string) error {
	if s == nil || s.client == nil || len(keys) == 0 {
		return nil
	}
	err := s.client.Del(ctx, keys...).Err()
	return NormalizeErr(err)
}

func (s *redisSession) Exists(ctx context.Context, key string) (bool, error) {
	if s == nil || s.client == nil {
		return false, nil
	}
	n, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, NormalizeErr(err)
	}
	return n > 0, nil
}

func (s *redisSession) Publish(ctx context.Context, channel string, message any) error {
	if s == nil || s.client == nil {
		return nil
	}
	err := s.client.Publish(ctx, channel, message).Err()
	return NormalizeErr(err)
}
