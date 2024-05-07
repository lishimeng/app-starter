package cms

import (
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/app-starter/tool"
)

// ConfigCategory 配置类型: file/db/redis/...
type ConfigCategory string

type OptionFunc func()

const (
	CategoryConfigFile = "config_file"
	CategoryDatabase   = "config_db"
	CategoryRedis      = "config_redis"
)

var category ConfigCategory

var name WebSiteName

// setName 设置app名称
func setName(module string) {
	if len(module) <= 0 {
		panic("cms name nil")
	}
	name = WebSiteName(module)
}

func WithName(name string) OptionFunc {
	return func() {
		setName(name)
	}
}

func WithConfigFile(c WebSiteInfo) OptionFunc {
	return func() {
		fromConfig(c)
		category = CategoryConfigFile
	}
}

func WithDatabase() OptionFunc {
	return func() {
		category = CategoryDatabase
	}
}

func WithRedis() OptionFunc {
	return func() {
		category = CategoryRedis
	}
}

func Init(opts ...OptionFunc) {
	for _, opt := range opts {
		opt()
	}
}

func GetWebsiteInfo() (ws WebSiteInfo, err error) {

	switch category {
	case CategoryConfigFile:
		ws = c
	case CategoryDatabase:
		ws, err = getWebsiteFromDb(name)
	case CategoryRedis:
	default:
		ws = getDefaultConfig()
	}

	return
}

// Router Theme提供的默认接口
func Router(p server.Router) {
	p.Get("/theme", ApiGetSpaConfig)
}

// ApiGetSpaConfig 开放接口：获取页面主题配置。
func ApiGetSpaConfig(ctx server.Context) {
	var resp SpaResp
	skipCache := ctx.C.URLParamIntDefault("skipCache", 0) //0-默认从缓存获取；1-跳过缓存；
	//默认AppName

	var themeConfigs []SpaConfigInfo
	switch skipCache {
	case 0:
		themeConfigs = GetPageTheme()
	default:
		themeConfigs = GetPageThemeSkipCache()
	}
	resp.Data = FormatPageTheme(themeConfigs)
	resp.Code = tool.RespCodeSuccess
	ctx.Json(resp)
}
