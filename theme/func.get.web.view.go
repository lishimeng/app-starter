package theme

import (
	"encoding/json"
	"fmt"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/app-starter/factory"
	"github.com/lishimeng/go-log"
)

func GetWebViewData() (data map[string]interface{}, err error) {
	if existWebViewCache() {
		data, err = getWebViewCache()
		if err == nil {
			return
		}
	}
	log.Debug(err)
	data, err = GetWebViewDataSkipCache()
	return
}

func GetWebViewDataSkipCache() (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	var wv AppWebView
	wv.AppName = AppName
	err = app.GetOrm().Context.Read(&wv, "AppName")
	if err != nil {
		return
	}
	if len(wv.ViewData) != 0 {
		bs, err := json.Marshal(wv.ViewData)
		if err != nil {
			return
		}
		err = json.Unmarshal(bs, &data)
		if err != nil {
			return
		}
	}
	return
}

func setWebViewCache(viewData map[string]interface{}) (err error) {
	if factory.GetCache() == nil {
		return
	}
	data, err := json.Marshal(viewData)
	if err != nil {
		return
	}
	return factory.GetCache().GetSkipLocal(webviewKey(), data)
}

func existWebViewCache() bool {
	if factory.GetCache() == nil {
		return false
	}
	return factory.GetCache().Exists(webviewKey())
}

func getWebViewCache() (viewData map[string]interface{}, err error) {
	if factory.GetCache() == nil {
		return
	}
	data := make([]byte, 0)
	err = factory.GetCache().GetSkipLocal(webviewKey(), &data)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &viewData)
	if err != nil {
		return
	}
	return
}

func webviewKey() string {
	return fmt.Sprintf(webViewCacheKeyTpl, AppName)
}
