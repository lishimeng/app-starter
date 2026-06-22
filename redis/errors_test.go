package redis

import (
	"errors"
	"testing"

	goredis "github.com/redis/go-redis/v9"
)

func TestIsRedisNil(t *testing.T) {
	if !IsRedisNil(goredis.Nil) {
		t.Fatal("expected redis.Nil")
	}
	if IsRedisNil(nil) || IsRedisNil(errors.New("boom")) {
		t.Fatal("expected false for nil and other errors")
	}
}
