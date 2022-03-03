package app

import (
	"context"
	shutdown "github.com/lishimeng/go-app-shutdown"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	var a = New()
	e := a.Start(func(ctx context.Context, builder *ApplicationBuilder) error {
		builder.EnableWeb(":8111")
		return nil
	}, func(s string) {
		t.Log(s)
	})
	if e != nil {
		t.Fatal(e)
	}

	go time.AfterFunc(time.Second*20, func() {
		shutdown.Exit("time to shutdown")
	})
}

func TestGetOrm(t *testing.T) {
	o := GetOrm()
	t.Logf("Orm is :%T", o)
}

func TestGetCache(t *testing.T) {
	o := GetCache()
	t.Logf("Cache is :%T", o)
}
