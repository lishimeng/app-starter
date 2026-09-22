package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lishimeng/app-starter/server/pprof"
)

type Handler func(ctx Context)

type Router interface {
	Get(path string, handlers ...Handler)
	Post(path string, handlers ...Handler)
	Put(path string, handlers ...Handler)
	Delete(path string, handlers ...Handler)
	Patch(path string, handlers ...Handler)
	Options(path string, handlers ...Handler)
	Any(path string, handlers ...Handler)
	Path(prefix string) Router
	MountPprof(relativePath string)
}

type router struct {
	g *gin.RouterGroup
}

func NewRouter(engine *gin.Engine) Router {
	return &router{g: engine.Group("")}
}

func (r *router) Path(prefix string) Router {
	return &router{g: r.g.Group(prefix)}
}

func (r *router) add(path string, handlers []Handler, register func(...gin.HandlerFunc)) {
	withRouteHandlerName(handlers, func() {
		register(wrapGinHandlers(handlers...)...)
	})
}

func (r *router) Get(path string, handlers ...Handler) {
	r.add(path, handlers, func(hs ...gin.HandlerFunc) { r.g.GET(path, hs...) })
}

func (r *router) Post(path string, handlers ...Handler) {
	r.add(path, handlers, func(hs ...gin.HandlerFunc) { r.g.POST(path, hs...) })
}

func (r *router) Put(path string, handlers ...Handler) {
	r.add(path, handlers, func(hs ...gin.HandlerFunc) { r.g.PUT(path, hs...) })
}

func (r *router) Delete(path string, handlers ...Handler) {
	r.add(path, handlers, func(hs ...gin.HandlerFunc) { r.g.DELETE(path, hs...) })
}

func (r *router) Patch(path string, handlers ...Handler) {
	r.add(path, handlers, func(hs ...gin.HandlerFunc) { r.g.PATCH(path, hs...) })
}

func (r *router) Options(path string, handlers ...Handler) {
	r.add(path, handlers, func(hs ...gin.HandlerFunc) { r.g.OPTIONS(path, hs...) })
}

func (r *router) Any(path string, handlers ...Handler) {
	r.add(path, handlers, func(hs ...gin.HandlerFunc) { r.g.Any(path, hs...) })
}

func (r *router) MountPprof(relativePath string) {
	pprof.Register(r.g, relativePath)
}
