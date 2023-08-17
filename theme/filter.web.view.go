package theme

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/go-log"
)

//WithWebViewFilter webViewFilter
func WithWebViewFilter(ctx iris.Context) {
	data, err := GetWebViewData()
	if err != nil {
		log.Debug(err)
		return
	}
	for k, v := range data {
		ctx.ViewData(k, v)
	}
	ctx.Next()
}
