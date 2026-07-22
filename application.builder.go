package app

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/lishimeng/app-starter/application/api"
	"github.com/lishimeng/app-starter/cache"
	"github.com/lishimeng/app-starter/log"
	"github.com/lishimeng/app-starter/mqtt"
	"github.com/lishimeng/app-starter/persistence"
	"github.com/lishimeng/app-starter/redis"
	"github.com/lishimeng/app-starter/server"
	"github.com/lishimeng/app-starter/token"
	"github.com/lishimeng/app-starter/version"
	"github.com/lishimeng/x/etc"
)

type TokenValidatorInjectFunc func(storage token.Storage)
type TokenValidatorBuilder func(injectFunc TokenValidatorInjectFunc)

// WithEmbed adapts embed.FS to EnableStaticWeb (files at FS root → URL /, /css/..., etc.).
func WithEmbed(fs embed.FS) func() http.FileSystem {
	return WithEmbedRoot(fs, ".")
}

// WithEmbedRoot serves a subdirectory inside fsys at web root.
// Use when embedding with //go:embed all:static and paths are static/index.html, etc.
func WithEmbedRoot(fsys fs.FS, root string) func() http.FileSystem {
	if root == "" || root == "." {
		return func() http.FileSystem { return http.FS(fsys) }
	}
	sub, err := fs.Sub(fsys, root)
	if err != nil {
		panic(fmt.Sprintf("app.WithEmbedRoot(%q): %v", root, err))
	}
	return func() http.FileSystem { return http.FS(sub) }
}

// EnableEmbedStaticWeb is shorthand for EnableStaticWeb(WithEmbedRoot(fs, root...)).
func (h *ApplicationBuilder) EnableEmbedStaticWeb(fsys fs.FS, root ...string) *ApplicationBuilder {
	dir := "."
	if len(root) > 0 && root[0] != "" {
		dir = root[0]
	}
	return h.EnableStaticWeb(WithEmbedRoot(fsys, dir))
}

type ApplicationBuilder struct {
	webEnable     bool
	webListen     string
	webComponents []server.Component

	webStaticEnable bool
	vdir            string
	assetInfo       func(name string) (os.FileInfo, error)
	asset           func(string) ([]byte, error)
	assetNames      func() []string
	webStaticHome   string
	assetFile       func() http.FileSystem

	webLogLevel string

	pprofListen         string
	pprofListenOverride bool
	adminSetup          server.AdminSetup
	stripTrailingSlash  bool

	dbEnable bool
	dbConfig persistence.BaseConfig
	dbModels []any
	dbViews  []any
	dbDebug  bool

	redisEnable bool
	redisOpts   redis.Options
	cacheEnable bool
	cacheOpts   cache.Options

	tokenValidatorEnable  bool
	tokenValidatorBuilder TokenValidatorBuilder

	mqttEnable  bool
	mqttOptions []mqtt.ClientOption

	// other components
	componentsBeforeWebServer []func(ctx context.Context) (err error)
	componentsAfterWebServer  []func(ctx context.Context) (err error)
}

var WithDefaultCallback = func(configName string) (f func(loader etc.Loader)) {
	return func(loader etc.Loader) {
		loader.SetFileSearcher(configName, ".").SetEnvPrefix("").SetEnvSearcher()
	}
}

func (h *ApplicationBuilder) LoadConfig(config interface{}, callback func(etc.Loader)) (err error) {
	var cb = callback
	loader := etc.New()
	if cb == nil {
		cb = WithDefaultCallback("config")
	}
	cb(loader)
	err = loader.Load(config)
	return
}

func (h *ApplicationBuilder) EnableWeb(listen string, components ...server.Component) *ApplicationBuilder {
	h.webEnable = true
	h.webListen = listen
	h.webComponents = append(h.webComponents, api.Router)
	if len(components) > 0 {
		h.webComponents = append(h.webComponents, components...)
	}
	// TODO check
	return h
}

func (h *ApplicationBuilder) SetMonitorPrefix(prefix string) *ApplicationBuilder {
	api.MonitorPrefix = prefix
	return h
}

func (h *ApplicationBuilder) HealthyHandler(handler func() int) *ApplicationBuilder {
	if handler != nil {
		api.LivenessHandler = handler
	}
	return h
}

func (h *ApplicationBuilder) ReadyHandler(handler func() int) *ApplicationBuilder {
	if handler != nil {
		api.ReadinessHandler = handler
	}
	return h
}

func (h *ApplicationBuilder) SetWebLogLevel(lvl string) *ApplicationBuilder {
	h.webLogLevel = lvl
	return h
}

// EnableStripTrailingSlash rewrites paths like /api/ → /api on web and admin
// before routing (nginx-style, no 301). Off by default.
func (h *ApplicationBuilder) EnableStripTrailingSlash() *ApplicationBuilder {
	h.stripTrailingSlash = true
	return h
}

// SetPprofListen overrides the default pprof listen address (DefaultPprofListen, :6060).
func (h *ApplicationBuilder) SetPprofListen(listen string) *ApplicationBuilder {
	h.pprofListen = listen
	h.pprofListenOverride = true
	return h
}

// EnableAdminRoutes registers custom routes on the admin listener (default :6060).
// Framework routes (/pprof, /metrics, POST /cl) are registered first; setup runs after.
// setup receives server.Router so business code does not import gin.
func (h *ApplicationBuilder) EnableAdminRoutes(setup server.AdminSetup) *ApplicationBuilder {
	h.adminSetup = setup
	return h
}

func (h *ApplicationBuilder) pprofListenAddr() string {
	if h.pprofListenOverride {
		return h.pprofListen
	}
	return server.DefaultPprofListen
}

func (h *ApplicationBuilder) EnableStaticWeb(assetFile func() http.FileSystem) *ApplicationBuilder {
	h.webStaticEnable = true
	h.assetFile = assetFile
	return h
}

func (h *ApplicationBuilder) EnableDatabase(config persistence.BaseConfig,
	models ...interface{}) *ApplicationBuilder {

	h.dbEnable = true
	h.dbConfig = config
	h.dbModels = models
	return h
}

func (h *ApplicationBuilder) EnableDatabaseView(views ...interface{}) *ApplicationBuilder {
	h.dbViews = views
	return h
}

// EnableDatabaseLog 鎵撳紑orm log杈撳嚭
func (h *ApplicationBuilder) EnableDatabaseLog() *ApplicationBuilder {
	h.dbDebug = true
	return h
}

func (h *ApplicationBuilder) EnableRedis(opts redis.Options) *ApplicationBuilder {
	h.redisEnable = true
	h.redisOpts = opts
	return h
}

func (h *ApplicationBuilder) EnableCache(cacheOpts cache.Options) *ApplicationBuilder {
	h.cacheEnable = true
	h.cacheOpts = cacheOpts
	return h
}

func (h *ApplicationBuilder) EnableOrmLog() *ApplicationBuilder {
	h.dbDebug = true
	return h
}

func (h *ApplicationBuilder) EnableMqtt(options ...mqtt.ClientOption) *ApplicationBuilder {
	h.mqttEnable = true
	if len(options) > 0 {
		h.mqttOptions = append(h.mqttOptions, options...)
	}
	log.Debug("enable mqtt module")
	return h
}

// EnableTokenValidator 楠岃瘉Token锛屼娇鐢≧edisTokenValidator鍓嶉渶瑕乪nableCache
func (h *ApplicationBuilder) EnableTokenValidator(builder TokenValidatorBuilder) *ApplicationBuilder {
	h.tokenValidatorEnable = true
	h.tokenValidatorBuilder = builder
	return h
}

func (h *ApplicationBuilder) PrintVersion() *ApplicationBuilder {
	version.Print()
	return h
}
