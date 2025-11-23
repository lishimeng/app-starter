package server

import "github.com/kataras/iris/v12"

type Handler func(ctx Context)

type Router interface {
	Post(path string, middleware ...Handler)
	Get(path string, middleware ...Handler)
	Path(path string) Router
	Party() iris.Party
	Put(path string, middleware ...Handler)
	Delete(path string, middleware ...Handler)
}

type router struct {
	p iris.Party
}

func (r *router) Party() iris.Party {
	return r.p
}

func (r *router) Post(path string, middleware ...Handler) {
	var handlers []iris.Handler
	for _, h := range middleware {
		handlers = append(handlers, func(ctx iris.Context) {
			h(Context{C: ctx})
		})
	}
	r.p.Post(path, handlers...)
}

func (r *router) Put(path string, middleware ...Handler) {
	var handlers []iris.Handler
	for _, h := range middleware {
		handlers = append(handlers, func(ctx iris.Context) {
			h(Context{C: ctx})
		})
	}
	r.p.Put(path, handlers...)
}

func (r *router) Delete(path string, middleware ...Handler) {
	var handlers []iris.Handler
	for _, h := range middleware {
		handlers = append(handlers, func(ctx iris.Context) {
			h(Context{C: ctx})
		})
	}
	r.p.Delete(path, handlers...)
}

func (r *router) Get(path string, middleware ...Handler) {
	var handlers []iris.Handler
	for _, h := range middleware {
		handlers = append(handlers, func(ctx iris.Context) {
			h(Context{C: ctx})
		})
	}
	r.p.Get(path, handlers...)
}

func (r *router) Patch(path string, middleware ...Handler) {
	var handlers []iris.Handler
	for _, h := range middleware {
		handlers = append(handlers, func(ctx iris.Context) {
			h(Context{C: ctx})
		})
	}
	r.p.Patch(path, handlers...)
}

func (r *router) Options(path string, middleware ...Handler) {
	var handlers []iris.Handler
	for _, h := range middleware {
		handlers = append(handlers, func(ctx iris.Context) {
			h(Context{C: ctx})
		})
	}
	r.p.Options(path, handlers...)
}

func NewRouter(app *iris.Application) Router {
	return &router{p: app.Party("/")}
}

func (r *router) Path(path string) Router {
	p := r.p.Party(path)
	return &router{p: p}
}
