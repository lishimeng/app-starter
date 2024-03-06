package api

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/server"
)

var ReadinessHandler func() int

func Ready(ctx server.Context) {
	if LivenessHandler != nil {
		ctx.C.StatusCode(ReadinessHandler())
	} else {
		ctx.C.StatusCode(iris.StatusOK)
	}
}
