package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestClient_GetSetDel(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	client := NewClient(Options{Addr: mr.Addr()})
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	if _, err := client.Get(ctx, "missing").Result(); !IsRedisNil(err) {
		t.Fatalf("expected not found, got %v", err)
	}

	if err := client.Set(ctx, "k", "v", time.Minute).Err(); err != nil {
		t.Fatal(err)
	}
	val, err := client.Get(ctx, "k").Result()
	if err != nil || val != "v" {
		t.Fatalf("got %q err=%v", val, err)
	}
	if err := client.Del(ctx, "k").Err(); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Get(ctx, "k").Result(); !IsRedisNil(err) {
		t.Fatalf("after del: %v", err)
	}
}

func TestClient_LPushRPop(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	client := NewClient(Options{Addr: mr.Addr()})
	defer func() { _ = client.Close() }()

	ctx := context.Background()
	if err := client.LPush(ctx, "q", "a", "b").Err(); err != nil {
		t.Fatal(err)
	}
	v, err := client.RPop(ctx, "q").Result()
	if err != nil || v != "a" {
		t.Fatalf("got %q err=%v", v, err)
	}
}
