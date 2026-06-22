package cache

import (
	"context"
	"time"
)

// RedisSession is the low-level Redis command surface (KV, pub/sub).
type RedisSession interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Publish(ctx context.Context, channel string, message any) error
}
