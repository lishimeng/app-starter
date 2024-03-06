package api

import (
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/lishimeng/app-starter/server"
)

var MonitorPrefix = "/m"

func Router(app server.Router) {
	p := app.Path(MonitorPrefix)
	//p := app.Party(MonitorPrefix)
	monitor(p)
}

func monitor(p server.Router) {
	p.Get("/healthy", Healthy)
	p.Get("/ready", Ready)
	p.Party().HandleMany("GET", "/debug/pprof /debug/pprof/{action:path}", pprof.New())
	p.Post("/cl", changeLogLevel)
}
