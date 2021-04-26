package cache

import (
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisOptions redis.Options

type Options struct {
	MaxSize int
	Ttl     time.Duration
}

const (
	defaultMaxSize = 10000
)

// New
/**
c := New(ctx, RedisOptions{
		Addr: addr,
		Password: "redisUAT",
	}, Options{
		MaxSize: 10,
		Ttl:     time.Second*20,
	})
 */
func New(ctx context.Context, redisOpts RedisOptions, cacheOpts Options) (c C) {
	var options = redis.Options(redisOpts)
	r := redis.NewClient(&options)

	if cacheOpts.MaxSize <= 0 {
		cacheOpts.MaxSize = defaultMaxSize
	}
	ca := cache.New(&cache.Options{
		Redis: r,
		LocalCache: cache.NewTinyLFU(cacheOpts.MaxSize, cacheOpts.Ttl),
	})
	rc := redisCache{
		ctx:        ctx,
		proxy:      ca,
		defaultTtl: cacheOpts.Ttl,
	}
	c = &rc
	return
}
