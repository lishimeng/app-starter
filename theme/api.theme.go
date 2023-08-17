package theme

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/go-log"
)

//Router Theme提供的默认接口
func Router(p iris.Party) {
	p.Get("/theme", ApiGetPageTheme)
}

//ApiGetPageTheme 开放接口：获取页面主题配置。
func ApiGetPageTheme(ctx iris.Context) {
	var resp response
	page := ctx.URLParamDefault("page", "")
	skipCache := ctx.URLParamIntDefault("skipCache", 0) //0-默认从缓存获取；1-跳过缓存；
	//默认AppName
	if len(AppName) == 0 || len(page) == 0 {
		log.Info("Not found app:%s, page:%s", AppName, page)
		resp.Code = tool.RespCodeNotFound
		resp.Message = fmt.Sprintf("Not found theme, app:%s, page:%s", AppName, page)
		tool.ResponseJSON(ctx, resp)
	}
	var themeConfigs []themeConfig
	switch skipCache {
	case 0:
		themeConfigs = GetPageTheme(page)
	default:
		themeConfigs = GetPageThemeSkipCache(page)
	}
	resp.Data = FormatPageTheme(themeConfigs)
	resp.Code = tool.RespCodeSuccess
	tool.ResponseJSON(ctx, resp)
}
