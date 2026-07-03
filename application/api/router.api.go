package api

import (
	"github.com/lishimeng/app-starter/server"
)

var MonitorPrefix = "/m"

func Router(app server.Router) {
	p := app.Path(MonitorPrefix)
	monitor(p)
}

func monitor(p server.Router) {
	p.Get("/healthy", Healthy)
	p.Get("/ready", Ready)
}
