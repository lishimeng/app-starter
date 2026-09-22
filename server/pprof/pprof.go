package pprof

import (
	httppprof "net/http/pprof"

	"github.com/gin-gonic/gin"
)

// Named gin adapters for net/http/pprof. Gin debug logs use runtime.FuncForPC
// on the last handler; gin.WrapF/WrapH are anonymous and print as WrapF.func1.

func Index(c *gin.Context)   { httppprof.Index(c.Writer, c.Request) }
func Cmdline(c *gin.Context) { httppprof.Cmdline(c.Writer, c.Request) }
func Profile(c *gin.Context) { httppprof.Profile(c.Writer, c.Request) }
func Symbol(c *gin.Context)  { httppprof.Symbol(c.Writer, c.Request) }
func Trace(c *gin.Context)   { httppprof.Trace(c.Writer, c.Request) }

func Allocs(c *gin.Context)       { httppprof.Handler("allocs").ServeHTTP(c.Writer, c.Request) }
func Block(c *gin.Context)        { httppprof.Handler("block").ServeHTTP(c.Writer, c.Request) }
func Goroutine(c *gin.Context)    { httppprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request) }
func Heap(c *gin.Context)         { httppprof.Handler("heap").ServeHTTP(c.Writer, c.Request) }
func Mutex(c *gin.Context)        { httppprof.Handler("mutex").ServeHTTP(c.Writer, c.Request) }
func Threadcreate(c *gin.Context) { httppprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request) }

// Register mounts standard pprof routes on r under prefix (e.g. "/pprof").
func Register(r gin.IRouter, prefix string) {
	g := r.Group(prefix)
	g.GET("/", Index)
	g.GET("/cmdline", Cmdline)
	g.GET("/profile", Profile)
	g.POST("/symbol", Symbol)
	g.GET("/symbol", Symbol)
	g.GET("/trace", Trace)
	g.GET("/allocs", Allocs)
	g.GET("/block", Block)
	g.GET("/goroutine", Goroutine)
	g.GET("/heap", Heap)
	g.GET("/mutex", Mutex)
	g.GET("/threadcreate", Threadcreate)
}
