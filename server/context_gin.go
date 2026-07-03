package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ginContext struct {
	*gin.Context
}

func (c *ginContext) BindJSON(v any) error {
	return c.ShouldBindJSON(v)
}

func (c *ginContext) QueryInt(name string) (int, error) {
	return strconv.Atoi(c.Query(name))
}

func (c *ginContext) ParamInt(name string) (int, error) {
	return strconv.Atoi(c.Param(name))
}

func (c *ginContext) JSON(resp any) {
	c.Context.JSON(http.StatusOK, resp)
}

func (c *ginContext) Redirect(url string, code int) {
	c.Context.Redirect(code, url)
}

func (c *ginContext) SetHeader(key, value string) {
	c.Header(key, value)
}

func (c *ginContext) Write(data []byte) (int, error) {
	return c.Writer.Write(data)
}

func (c *ginContext) GetCookie(name string) (string, error) {
	return c.Context.Cookie(name)
}

func (c *ginContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Writer, cookie)
}

func (c *ginContext) FormValue(name string) string {
	return c.Context.Request.FormValue(name)
}

func (c *ginContext) Request() *http.Request {
	return c.Context.Request
}

func (c *ginContext) ResponseWriter() http.ResponseWriter {
	return c.Writer
}

func (c *ginContext) Path() string {
	return c.Context.Request.URL.Path
}

func (c *ginContext) Html(layout, view string, data any) error {
	if layout != "" {
		return fmt.Errorf("server: layout %q not supported with gin templates", layout)
	}
	c.HTML(http.StatusOK, view, data)
	return nil
}

// GinHandler adapts a server.Handler to gin.HandlerFunc (e.g. admin listener routes).
func GinHandler(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(wrapContext(&ginContext{Context: c}))
	}
}

func wrapGinHandler(h Handler) gin.HandlerFunc {
	return GinHandler(h)
}

func wrapGinHandlers(handlers ...Handler) []gin.HandlerFunc {
	out := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		out = append(out, wrapGinHandler(h))
	}
	return out
}
