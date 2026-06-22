package cache

import (
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNormalizeErr_RedisNil(t *testing.T) {
	got := NormalizeErr(redis.Nil)
	if !IsNotFound(got) {
		t.Fatalf("expected not found, got %v", got)
	}
	if IsRedisNil(got) {
		t.Fatal("normalized err must not be redis.Nil")
	}
}

func TestNormalizeErr_Nil(t *testing.T) {
	if NormalizeErr(nil) != nil {
		t.Fatal("expected nil")
	}
}

func TestIsRedisNil(t *testing.T) {
	if !IsRedisNil(redis.Nil) {
		t.Fatal("expected redis.Nil")
	}
	if IsRedisNil(NormalizeErr(redis.Nil)) {
		t.Fatal("normalized err must not match redis.Nil")
	}
	if IsRedisNil(nil) || IsRedisNil(errors.New("boom")) {
		t.Fatal("expected false for nil and other errors")
	}
}

func TestIsNotFound(t *testing.T) {
	normalized := NormalizeErr(redis.Nil)
	if !IsNotFound(normalized) {
		t.Fatal("expected normalized not found")
	}
	if IsNotFound(redis.Nil) {
		t.Fatal("raw redis.Nil is not cache-level ErrNotFound")
	}
	if !errors.Is(normalized, ErrNotFound) {
		t.Fatal("normalized err must be ErrNotFound")
	}
}

func TestIsNotFoundAny(t *testing.T) {
	if !IsNotFoundAny(redis.Nil) {
		t.Fatal("expected raw redis.Nil")
	}
	if !IsNotFoundAny(NormalizeErr(redis.Nil)) {
		t.Fatal("expected normalized not found")
	}
	if IsNotFoundAny(nil) || IsNotFoundAny(errors.New("boom")) {
		t.Fatal("expected false for nil and other errors")
	}
}
