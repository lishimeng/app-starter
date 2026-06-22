package cache

import (
	"errors"
	"testing"

	goredis "github.com/redis/go-redis/v9"
)

func TestNormalizeErr_RedisNil(t *testing.T) {
	got := NormalizeErr(goredis.Nil)
	if !IsNotFound(got) {
		t.Fatalf("expected not found, got %v", got)
	}
}

func TestNormalizeErr_Nil(t *testing.T) {
	if NormalizeErr(nil) != nil {
		t.Fatal("expected nil")
	}
}

func TestIsNotFound(t *testing.T) {
	normalized := NormalizeErr(goredis.Nil)
	if !IsNotFound(normalized) {
		t.Fatal("expected normalized not found")
	}
	if IsNotFound(goredis.Nil) {
		t.Fatal("raw redis.Nil is not cache-level ErrNotFound")
	}
	if !errors.Is(normalized, ErrNotFound) {
		t.Fatal("normalized err must be ErrNotFound")
	}
}

func TestIsNotFoundAny(t *testing.T) {
	if !IsNotFoundAny(goredis.Nil) {
		t.Fatal("expected raw redis.Nil")
	}
	if !IsNotFoundAny(NormalizeErr(goredis.Nil)) {
		t.Fatal("expected normalized not found")
	}
	if IsNotFoundAny(nil) || IsNotFoundAny(errors.New("boom")) {
		t.Fatal("expected false for nil and other errors")
	}
}
