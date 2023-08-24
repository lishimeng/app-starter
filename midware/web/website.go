package web

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/cms"
)

var Name cms.WebSiteName

const Frame = "WebsiteFrame"

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
	ctx.ViewData(Frame, ws)
	ctx.Next()
}
