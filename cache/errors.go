package cache

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var errNotFound = errors.New("cache: not found")

// NormalizeErr maps redis driver errors to cache-level errors.
func NormalizeErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, redis.Nil) {
		return errNotFound
	}
	return err
}

// IsNotFound reports whether err is a missing-key error from cache/redis session.
func IsNotFound(err error) bool {
	return errors.Is(err, errNotFound)
}
