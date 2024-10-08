package server

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/go-log"
	"net/http"
	"time"
)

const (
	defaultMonitorAddr = ":8888"
)

//type Component func(app *iris.Application)

type Component func(root Router)
type Config struct {
	Listen string
	LogLvl string
}

type Server struct {
	config  Config
	proxy   *iris.Application
	monitor *iris.Application
}

func New(config Config) (handler *Server) {

	s := Server{
		config:  config,
		proxy:   iris.New(),
		monitor: iris.New(),
	}
	s.setLogLvl(config.LogLvl)
	return &s
}

// GetApplication 服务器实例
func (s *Server) GetApplication() *iris.Application {
	return s.proxy
}

// GetMonitor 监控实例
func (s *Server) GetMonitor() *iris.Application {
	return s.monitor
}

func (s *Server) setLogLvl(lvl string) *Server {
	s.proxy.Logger().SetLevel(lvl)
	// Monitor不设置log
	return s
}

func (s *Server) RegisterComponent(component Component) *Server {
	var r = NewRouter(s.proxy)
	component(r)
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
	if err := s.monitor.Configure(iris.WithCharset("UTF-8")).Build(); err != nil {
		return err
	}
	monitorServer := http.Server{Addr: defaultMonitorAddr, Handler: s.monitor}
	go s.shutdownFuture(&srv, ctx)
	go s.shutdownFuture(&monitorServer, ctx)

	go func() {
		_ = monitorServer.ListenAndServe()
	}()

	log.Info("web server listen %s", s.config.Listen)
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
			c = context.TODO()
			if err := srv.Shutdown(c); nil != err {
			}
			return
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}
