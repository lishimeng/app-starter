package server

import (
	"context"
	"github.com/kataras/iris/v12"
	"net/http"
	"time"
)

type Component func(app *iris.Application)
type Config struct {
	Listen string
	LogLvl string
}

type Server struct {
	config        Config
	proxy         *iris.Application
}

func New(config Config) (handler *Server) {

	s := Server{
		config: config,
		proxy:  iris.New(),
	}
	s.setLogLvl(config.LogLvl)
	return &s
}

func (s *Server) setLogLvl(lvl string) *Server {
	s.proxy.Logger().SetLevel(lvl)
	return s
}

func (s *Server) RegisterComponent(component Component) *Server {
	component(s.proxy)
	return s
}

func (s *Server) AdvancedConfig(handler func(app *iris.Application)) *Server {

	if handler != nil {
		handler(s.proxy)
	}
	return s
}

func (s *Server) SetHomePage(indexHtml string) *Server {
	s.proxy.Get("/", func(c iris.Context) {
		_, _ = c.HTML(indexHtml)
	})
	return s
}

func (s *Server) OnErrorCode(code int, onErr func(ctx iris.Context)) *Server {
	s.proxy.OnErrorCode(code, onErr)
	return s
}

func (s *Server) RegisterComponents(components ...Component) *Server {

	if len(components) > 0 {
		for _, component := range components {
			s.RegisterComponent(component)
		}
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.proxy.Configure(iris.WithCharset("UTF-8")).Build(); err != nil {
		return err
	}
	srv := http.Server{
		Addr:    s.config.Listen,
		Handler: s.proxy,
	}
	go s.shutdownFuture(&srv, ctx)

	return srv.ListenAndServe()
}

func (s *Server) shutdownFuture(srv *http.Server, ctx context.Context) {
	if ctx == nil {
		return
	}
	var c context.Context
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			c, cancel = context.WithTimeout(context.Background(), 3*time.Second)
			if err := srv.Shutdown(c); nil != err {
			}
			return
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}
