package server

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

func init() {
	prev := gin.DebugPrintRouteFunc
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		if name := currentRouteHandlerName(); name != "" && isAdapterHandlerName(handlerName) {
			handlerName = name
		}
		if prev != nil {
			prev(httpMethod, absolutePath, handlerName, nuHandlers)
			return
		}
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] %-6s %-25s --> %s (%d handlers)\n",
			httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

var (
	routeNameMu      sync.Mutex
	routeHandlerName atomic.Value // string
)

func withRouteHandlerName(handlers []Handler, fn func()) {
	name := ""
	if n := len(handlers); n > 0 {
		name = nameOfFunction(handlers[n-1])
	}
	routeNameMu.Lock()
	defer routeNameMu.Unlock()
	routeHandlerName.Store(name)
	defer routeHandlerName.Store("")
	fn()
}

func currentRouteHandlerName() string {
	v, _ := routeHandlerName.Load().(string)
	return v
}

func isAdapterHandlerName(handlerName string) bool {
	return strings.Contains(handlerName, "GinHandler.func1")
}

func nameOfFunction(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
