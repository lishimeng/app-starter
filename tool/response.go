package tool

import "github.com/kataras/iris/v12"

func ResponseJSON(ctx iris.Context, j interface{}) {
	_ = ctx.JSON(j)
}
