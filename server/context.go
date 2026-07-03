package server

import (
	"mime/multipart"
	"net/http"
)

// Context is the framework-agnostic HTTP request context passed to handlers and middleware.
type Context struct {
	gin *ginContext
}

func wrapContext(c *ginContext) Context {
	return Context{gin: c}
}

// JSON writes resp as JSON with status 200.
func (c *Context) JSON(resp any) {
	if c == nil || c.gin == nil || resp == nil {
		return
	}
	c.gin.JSON(resp)
}

// Json is an alias for JSON.
func (c *Context) Json(resp any) {
	c.JSON(resp)
}

// BindJSON decodes the request body JSON into v.
func (c *Context) BindJSON(v any) error {
	if c == nil || c.gin == nil {
		return nil
	}
	return c.gin.BindJSON(v)
}

// Query returns the URL query parameter named name.
func (c *Context) Query(name string) string {
	if c == nil || c.gin == nil {
		return ""
	}
	return c.gin.Query(name)
}

// QueryInt parses the URL query parameter named name as an int.
func (c *Context) QueryInt(name string) (int, error) {
	if c == nil || c.gin == nil {
		return 0, nil
	}
	return c.gin.QueryInt(name)
}

// Param returns the path parameter named name (e.g. ":id" in the route).
func (c *Context) Param(name string) string {
	if c == nil || c.gin == nil {
		return ""
	}
	return c.gin.Param(name)
}

// ParamInt parses the path parameter named name as an int.
func (c *Context) ParamInt(name string) (int, error) {
	if c == nil || c.gin == nil {
		return 0, nil
	}
	return c.gin.ParamInt(name)
}

// FormFile returns the first file for the given multipart form key.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c == nil || c.gin == nil {
		return nil, http.ErrNotSupported
	}
	return c.gin.FormFile(name)
}

// FormValue returns the first value for the given form key.
func (c *Context) FormValue(name string) string {
	if c == nil || c.gin == nil {
		return ""
	}
	return c.gin.FormValue(name)
}

// Status sets the HTTP response status code without writing a body.
func (c *Context) Status(code int) {
	if c == nil || c.gin == nil {
		return
	}
	c.gin.Status(code)
}

// Redirect sends an HTTP redirect to url with the given status code.
func (c *Context) Redirect(url string, code int) {
	if c == nil || c.gin == nil {
		return
	}
	c.gin.Redirect(url, code)
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	if c == nil || c.gin == nil {
		return
	}
	c.gin.SetHeader(key, value)
}

// Write writes raw bytes to the response body.
func (c *Context) Write(data []byte) (int, error) {
	if c == nil || c.gin == nil {
		return 0, http.ErrNotSupported
	}
	return c.gin.Write(data)
}

// Next invokes the remaining handlers in the middleware chain.
func (c *Context) Next() {
	if c == nil || c.gin == nil {
		return
	}
	c.gin.Next()
}

// Set stores a value in the per-request context (visible to later middleware/handlers).
func (c *Context) Set(key string, value any) {
	if c == nil || c.gin == nil {
		return
	}
	c.gin.Set(key, value)
}

// Get retrieves a value previously stored with Set.
func (c *Context) Get(key string) (any, bool) {
	if c == nil || c.gin == nil {
		return nil, false
	}
	return c.gin.Get(key)
}

// GetHeader returns the request header named name.
func (c *Context) GetHeader(name string) string {
	if c == nil || c.gin == nil {
		return ""
	}
	return c.gin.GetHeader(name)
}

// ClientIP returns the client IP address of the request.
func (c *Context) ClientIP() string {
	if c == nil || c.gin == nil {
		return ""
	}
	return c.gin.ClientIP()
}

// Request returns the underlying HTTP request.
func (c *Context) Request() *http.Request {
	if c == nil || c.gin == nil {
		return nil
	}
	return c.gin.Request()
}

// ResponseWriter returns the underlying HTTP response writer.
func (c *Context) ResponseWriter() http.ResponseWriter {
	if c == nil || c.gin == nil {
		return nil
	}
	return c.gin.ResponseWriter()
}

// GetCookie returns the named cookie from the request.
func (c *Context) GetCookie(name string) (string, error) {
	if c == nil || c.gin == nil {
		return "", http.ErrNotSupported
	}
	return c.gin.GetCookie(name)
}

// SetCookie adds a Set-Cookie header to the response.
func (c *Context) SetCookie(cookie *http.Cookie) {
	if c == nil || c.gin == nil || cookie == nil {
		return
	}
	c.gin.SetCookie(cookie)
}

// Path returns the request URL path (without query string).
func (c *Context) Path() string {
	if c == nil || c.gin == nil {
		return ""
	}
	return c.gin.Path()
}

// Html renders a server-side template named view with data.
// layout is ignored; templates must be loaded via Server.LoadHTMLGlob.
func (c *Context) Html(layout, view string, data any) error {
	if c == nil || c.gin == nil {
		return nil
	}
	return c.gin.Html(layout, view, data)
}
