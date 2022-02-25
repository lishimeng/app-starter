package api

import (
	"github.com/kataras/iris/v12"
)

var LivenessHandler func() int

func Healthy(ctx iris.Context) {
	if LivenessHandler != nil {
		ctx.StatusCode(LivenessHandler())
	} else {
		ctx.StatusCode(iris.StatusOK)
	}
}
