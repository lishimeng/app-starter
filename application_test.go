package app

import (
	"context"
	shutdown "github.com/lishimeng/go-app-shutdown"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	go time.AfterFunc(time.Second*20, func() {
		shutdown.Exit("time to shutdown")
	})
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
	go time.AfterFunc(time.Second*20, func() {
		shutdown.Exit("time to shutdown")
	})
	var a = New()
	e := a.Start(func(ctx context.Context, builder *ApplicationBuilder) error {
		builder.EnableWeb(":8111").
			ComponentAfter(func(ctx context.Context) (err error) {
				proxy := GetWebServer().GetApplication()
				if proxy == nil {
					t.Fatal("web server nil")
					return
				} else {
					t.Logf("print web server:%s", proxy.String())
				}
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

func TestGetAmqp(t *testing.T) {
	o := GetAmqp()
	t.Logf("amqp is %+v", o)
}

func TestGetNamedOrm(t *testing.T) {
	o := GetNamedOrm("default")
	t.Logf("named orm is %+v", o)
}
