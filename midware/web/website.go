package web

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/cms"
	"github.com/lishimeng/app-starter/tool"
)

var Name cms.WebSiteName

func Website(ctx iris.Context) {

	if len(Name) <= 0 {
		ctx.Next()
		return
	}
	ws, err := cms.GetWebsiteFrame(Name)
	if err != nil {
		ctx.Next()
		return
	}
	ctx.ViewData(tool.WebsiteCtx, ws)
	ctx.Next()
}
