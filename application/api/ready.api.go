package api

import (
	"net/http"

	"github.com/lishimeng/app-starter/server"
)

var ReadinessHandler func() int

func Ready(ctx server.Context) {
	if ReadinessHandler != nil {
		ctx.Status(ReadinessHandler())
	} else {
		ctx.Status(http.StatusOK)
	}
}
