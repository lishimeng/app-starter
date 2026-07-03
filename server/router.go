package server

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
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

func (r *router) Get(path string, handlers ...Handler) {
	r.g.GET(path, wrapGinHandlers(handlers...)...)
}

func (r *router) Post(path string, handlers ...Handler) {
	r.g.POST(path, wrapGinHandlers(handlers...)...)
}

func (r *router) Put(path string, handlers ...Handler) {
	r.g.PUT(path, wrapGinHandlers(handlers...)...)
}

func (r *router) Delete(path string, handlers ...Handler) {
	r.g.DELETE(path, wrapGinHandlers(handlers...)...)
}

func (r *router) Patch(path string, handlers ...Handler) {
	r.g.PATCH(path, wrapGinHandlers(handlers...)...)
}

func (r *router) Options(path string, handlers ...Handler) {
	r.g.OPTIONS(path, wrapGinHandlers(handlers...)...)
}

func (r *router) Any(path string, handlers ...Handler) {
	r.g.Any(path, wrapGinHandlers(handlers...)...)
}

func (r *router) MountPprof(relativePath string) {
	pprof.Register(r.g, relativePath)
}
