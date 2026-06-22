package cache

import (
	"context"
	"time"
)

// RedisContext is the facade for RedisSession, similar to persistence.OrmContext.
type RedisContext struct {
	session RedisSession
}

func NewRedis() *RedisContext {
	if s := resolveRedisSession(); s != nil {
		return WrapRedisSession(s)
	}
	return &RedisContext{}
}

func WrapRedisSession(s RedisSession) *RedisContext {
	return &RedisContext{session: s}
}

func (r *RedisContext) Get(ctx context.Context, key string) (string, error) {
	if r == nil || r.session == nil {
		return "", nil
	}
	return r.session.Get(ctx, key)
}

func (r *RedisContext) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if r == nil || r.session == nil {
		return nil
	}
	return r.session.Set(ctx, key, value, ttl)
}

func (r *RedisContext) Del(ctx context.Context, keys ...string) error {
	if r == nil || r.session == nil {
		return nil
	}
	return r.session.Del(ctx, keys...)
}

func (r *RedisContext) Exists(ctx context.Context, key string) (bool, error) {
	if r == nil || r.session == nil {
		return false, nil
	}
	return r.session.Exists(ctx, key)
}

func (r *RedisContext) Publish(ctx context.Context, channel string, message any) error {
	if r == nil || r.session == nil {
		return nil
	}
	return r.session.Publish(ctx, channel, message)
}

var redisSessionResolver func() RedisSession

// SetRedisSessionResolver installs the global session resolver (used by factory).
func SetRedisSessionResolver(fn func() RedisSession) {
	redisSessionResolver = fn
}

func resolveRedisSession() RedisSession {
	if redisSessionResolver == nil {
		return nil
	}
	return redisSessionResolver()
}
