package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRedisSession_GetSetDel(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer func() {
		_ = client.Close()
	}()

	s := NewRedisSession(client)
	ctx := context.Background()

	if _, err := s.Get(ctx, "missing"); !IsNotFound(err) {
		t.Fatalf("expected not found, got %v", err)
	}

	if err := s.Set(ctx, "k", "v", time.Minute); err != nil {
		t.Fatal(err)
	}
	val, err := s.Get(ctx, "k")
	if err != nil || val != "v" {
		t.Fatalf("got %q err=%v", val, err)
	}
	ok, err := s.Exists(ctx, "k")
	if err != nil || !ok {
		t.Fatalf("exists: ok=%v err=%v", ok, err)
	}
	if err := s.Del(ctx, "k"); err != nil {
		t.Fatal(err)
	}
	ok, err = s.Exists(ctx, "k")
	if err != nil || ok {
		t.Fatalf("after del: ok=%v err=%v", ok, err)
	}
}

func TestRedisSession_Publish(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer func() {
		_ = client.Close()
	}()

	s := NewRedisSession(client)
	if err := s.Publish(context.Background(), "ch", "msg"); err != nil {
		t.Fatal(err)
	}
}

func TestRedisContext_NewRedis(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	SetRedisSessionResolver(func() RedisSession { return NewRedisSession(client) })

	r := NewRedis()
	if err := r.Set(context.Background(), "a", "1", time.Minute); err != nil {
		t.Fatal(err)
	}
	v, err := r.Get(context.Background(), "a")
	if err != nil || v != "1" {
		t.Fatalf("got %q err=%v", v, err)
	}
}
