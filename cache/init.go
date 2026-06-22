package cache

import (
	"context"
	"time"

	gocache "github.com/go-redis/cache/v9"
	starterredis "github.com/lishimeng/app-starter/redis"
)

type Options struct {
	MaxSize int
	Ttl     time.Duration
}

const defaultMaxSize = 10000

// New builds a cache.C from an existing redis.Client (shared with app.GetRedisClient).
func New(ctx context.Context, client *starterredis.Client, cacheOpts Options) (c C) {
	if client == nil || client.Client == nil {
		return nil
	}
	if cacheOpts.MaxSize <= 0 {
		cacheOpts.MaxSize = defaultMaxSize
	}
	ca := gocache.New(&gocache.Options{
		Redis:      client.Client,
		LocalCache: gocache.NewTinyLFU(cacheOpts.MaxSize, cacheOpts.Ttl),
	})
	rc := redisCache{
		ctx:        ctx,
		proxy:      ca,
		defaultTtl: cacheOpts.Ttl,
	}
	c = &rc
	return
}
