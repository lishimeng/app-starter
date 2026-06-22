package cache

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var errNotFound = errors.New("cache: not found")

// ErrNotFound is returned by framework APIs after NormalizeErr (not redis.Nil).
var ErrNotFound = errNotFound

// NormalizeErr maps redis driver errors to cache-level errors.
func NormalizeErr(err error) error {
	if err == nil {
		return nil
	}
	if IsRedisNil(err) {
		return errNotFound
	}
	return err
}

// IsRedisNil reports raw go-redis missing-key errors. Use for direct go-redis calls only.
func IsRedisNil(err error) bool {
	return err != nil && errors.Is(err, redis.Nil)
}

// IsNotFound reports normalized not-found errors from RedisContext / cache.C.
func IsNotFound(err error) bool {
	return err != nil && errors.Is(err, ErrNotFound)
}

// IsNotFoundAny reports not-found when err source is unknown or mixed.
func IsNotFoundAny(err error) bool {
	return IsRedisNil(err) || IsNotFound(err)
}
