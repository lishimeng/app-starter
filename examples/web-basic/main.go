package main

import (
	"context"
	"time"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/examples/model"
	"github.com/lishimeng/app-starter/examples/web-basic/proc"
	"github.com/lishimeng/app-starter/examples/web-basic/router"
	"github.com/lishimeng/app-starter/examples/web-basic/setup"
	"github.com/lishimeng/app-starter/log"
)

func main() {

	defer func() {
		time.Sleep(time.Millisecond * 200)
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%v", err)
		}
	}()

	err := _main()
	if err != nil {
		log.Errorf("%v", err)
	}

}

func _main() (err error) {
	log.Config().LevelFromString("INFO").JSON().Apply()

	application := app.New()

	err = application.Start(func(ctx context.Context, builder *app.ApplicationBuilder) (e error) {

		log.Info("application start")
		builder.
			EnableDatabase(setup.PostgresConfig().Build(), new(model.BusinessConnector)).
			EnableDatabaseLog().
			SetWebLogLevel("DEBUG").
			ComponentBefore(proc.Before).
			ComponentAfter(proc.After).
			EnableWeb(setup.WebPort(), router.Router).
			PrintVersion()
		return
	}, func(s string) {
		log.Info(s)
	})

	return
}
