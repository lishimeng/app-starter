package tool

import "github.com/kataras/iris/v12"

const (
	RespCodeSuccess  = 200
	RespCodeNotFound = 404
	RespCodeError    = 500
)

const (
	RespMsgNotFount = "not found"
	RespMsgIdNum    = "id must be a int value"
)

func ResponseJSON(ctx iris.Context, j interface{}) {
	_ = ctx.JSON(j)
}
