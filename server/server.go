package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Component func(root Router)

type Config struct {
	Listen             string
	LogLvl             string
	StripTrailingSlash bool // rewrite /api/ → /api before routing (no 301)
}

type Server struct {
	config Config
	engine *gin.Engine
}

func New(config Config) *Server {
	if strings.EqualFold(config.LogLvl, "debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Server{
		config: config,
		engine: gin.New(),
	}
	if config.StripTrailingSlash {
		s.engine.RedirectTrailingSlash = false
	}
	s.engine.Use(gin.Recovery())
	return s
}

func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

func (s *Server) RegisterComponent(component Component) *Server {
	component(NewRouter(s.engine))
	return s
}

func (s *Server) AdvancedConfig(handler func(*gin.Engine)) *Server {
	if handler != nil {
		handler(s.engine)
	}
	return s
}

func (s *Server) LoadHTMLGlob(pattern string) *Server {
	s.engine.LoadHTMLGlob(pattern)
	return s
}

func (s *Server) StaticFS(relativePath string, fs http.FileSystem) *Server {
	s.engine.StaticFS(relativePath, fs)
	return s
}

func (s *Server) SetNoRoute(handler Handler) *Server {
	s.engine.NoRoute(wrapGinHandler(handler))
	return s
}

// EnableMetrics installs HTTP RED Prometheus middleware on the main web server.
func (s *Server) EnableMetrics() *Server {
	s.engine.Use(MetricsMiddleware())
	return s
}

func (s *Server) SetHomePage(indexHTML string) *Server {
	s.engine.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHTML))
	})
	return s
}

func (s *Server) RegisterComponents(components ...Component) *Server {
	for _, component := range components {
		s.RegisterComponent(component)
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	handler := http.Handler(s.engine)
	if s.config.StripTrailingSlash {
		handler = stripTrailingSlash(handler)
	}
	srv := &http.Server{
		Addr:    s.config.Listen,
		Handler: handler,
	}
	return serveHTTP(ctx, srv, "web server")
}
