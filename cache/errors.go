package cache

import (
	"errors"

	starterredis "github.com/lishimeng/app-starter/redis"
)

var errNotFound = errors.New("cache: not found")

// ErrNotFound is returned by cache.C APIs after NormalizeErr (not redis.Nil).
var ErrNotFound = errNotFound

// NormalizeErr maps redis driver errors to cache-level errors.
func NormalizeErr(err error) error {
	if err == nil {
		return nil
	}
	if starterredis.IsRedisNil(err) {
		return errNotFound
	}
	return err
}

// IsNotFound reports normalized not-found errors from cache.C.
func IsNotFound(err error) bool {
	return err != nil && errors.Is(err, ErrNotFound)
}

// IsNotFoundAny reports not-found for raw redis.Nil or normalized cache errors.
func IsNotFoundAny(err error) bool {
	return starterredis.IsRedisNil(err) || IsNotFound(err)
}
