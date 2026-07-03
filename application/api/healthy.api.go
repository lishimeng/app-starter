package api

import (
	"net/http"

	"github.com/lishimeng/app-starter/server"
)

var LivenessHandler func() int

func Healthy(ctx server.Context) {
	if LivenessHandler != nil {
		ctx.Status(LivenessHandler())
	} else {
		ctx.Status(http.StatusOK)
	}
}
