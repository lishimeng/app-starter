package api

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/server"
)

var LivenessHandler func() int

func Healthy(ctx server.Context) {
	if LivenessHandler != nil {
		ctx.C.StatusCode(LivenessHandler())
	} else {
		ctx.C.StatusCode(iris.StatusOK)
	}
}
