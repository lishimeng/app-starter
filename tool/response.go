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

const WebsiteCtx = "WebSiteCtx"
const PageCtx = "PageCtx"

func ResponseJSON(ctx iris.Context, j interface{}) {
	_ = ctx.JSON(j)
}

func ResponseHtml(ctx iris.Context, layout string, view string, data any) {

	if len(layout) > 0 {
		ctx.ViewLayout(layout)
	}

	ctx.ViewData(PageCtx, data)
	_ = ctx.View(view)
}
