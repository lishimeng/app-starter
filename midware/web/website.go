package web

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/cms"
	"github.com/lishimeng/app-starter/tool"
)

// Website 注入网站的基本配置
func Website(opts ...cms.OptionFunc) func(iris.Context) {
	cms.Init(opts...)
	return func(ctx iris.Context) {
		ws, err := cms.GetWebsiteInfo()
		if err != nil {
			ctx.Next()
			return
		}
		ctx.ViewData(tool.WebsiteCtx, ws)
		ctx.Next()
	}

}
