package main

import (
	"context"
	"fmt"

	"time"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/examples/web-basic/proc"
	"github.com/lishimeng/app-starter/examples/web-basic/router"
	"github.com/lishimeng/go-log"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	err := _main()
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Millisecond * 200)
}

func _main() (err error) {
	application := app.New()

	err = application.Start(func(ctx context.Context, builder *app.ApplicationBuilder) (e error) {

		builder.
			SetWebLogLevel("DEBUG").
			ComponentBefore(proc.Before).
			ComponentAfter(proc.After).
			EnableWeb(":9527", router.Router).
			PrintVersion()
		return
	}, func(s string) {
		log.Info(s)
	})

	return
}
