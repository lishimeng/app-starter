package cache

import (
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNormalizeErr_RedisNil(t *testing.T) {
	got := NormalizeErr(redis.Nil)
	if !IsNotFound(got) {
		t.Fatalf("expected not found, got %v", got)
	}
}

func TestNormalizeErr_Nil(t *testing.T) {
	if NormalizeErr(nil) != nil {
		t.Fatal("expected nil")
	}
}
