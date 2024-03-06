package server

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/tool"
)

type Context struct {
	C iris.Context
}

func (c *Context) Json(resp any) {
	if resp == nil {
		return
	}
	_ = c.C.JSON(resp)
}

func (c *Context) Html(ctx Context, layout string, view string, data any) {

	if len(layout) > 0 {
		ctx.C.ViewLayout(layout)
	} else {
		ctx.C.ViewLayout(iris.NoLayout)
	}

	ctx.C.ViewData(tool.PageCtx, data)
	_ = ctx.C.View(view)
}
