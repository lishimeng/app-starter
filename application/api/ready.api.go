package api

import (
	"github.com/kataras/iris/v12"
)

var ReadinessHandler func() int

func Ready(ctx iris.Context) {
	if LivenessHandler != nil {
		ctx.StatusCode(ReadinessHandler())
	} else {
		ctx.StatusCode(iris.StatusOK)
	}
}
