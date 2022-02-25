package api

import "github.com/kataras/iris/v12"

var MonitorPrefix = "/m"

func Router(app *iris.Application) {
	p := app.Party(MonitorPrefix)
	monitor(p)
}

func monitor(p iris.Party) {
	p.Get("/healthy", Healthy)
	p.Get("/ready", Ready)
}
