package api

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
)

var MonitorPrefix = "/m"

func Router(app *iris.Application) {
	p := app.Party(MonitorPrefix)
	monitor(p)

}

func monitor(p iris.Party) {
	p.Get("/healthy", Healthy)
	p.Get("/ready", Ready)
	p.HandleMany("GET", "/debug/pprof /debug/pprof/{action:path}", pprof.New())
	p.Post("/cl", changeLogLevel)
}
