package app

import (
	"context"
	"testing"
	"time"

	shutdown "github.com/lishimeng/go-app-shutdown"
)

func exitAfterStart() {
	go func() {
		time.Sleep(100 * time.Millisecond)
		shutdown.Exit("test done")
	}()
}

func TestNew(t *testing.T) {
	exitAfterStart()
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
}

func TestGetWebServer(t *testing.T) {
	exitAfterStart()
	var a = New()
	e := a.Start(func(ctx context.Context, builder *ApplicationBuilder) error {
		builder.EnableWeb(":8111").
			ComponentAfter(func(ctx context.Context) (err error) {
				engine := GetWebServer().GetEngine()
				if engine == nil {
					t.Fatal("web server nil")
					return
				}
				t.Logf("web server engine ready")
				return
			})
		return nil
	}, func(s string) {
		t.Log(s)
	})
	if e != nil {
		t.Fatal(e)
	}
}

func TestGetNamedOrm(t *testing.T) {
	o := GetNamedOrm("default")
	t.Logf("named orm is %+v", o)
}
