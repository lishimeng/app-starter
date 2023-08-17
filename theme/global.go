package theme

var AppName string

const (
	pageThemeCacheKeyTpl = "theme_app-%s_page-%s"
	webViewCacheKeyTpl   = "theme_app-%s"
)

type response struct {
	Code    interface{} `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type themeConfig struct {
	Id                int         `json:"id,omitempty"`
	AppName           string      `json:"appName,omitempty"`
	ConfigPage        string      `json:"configPage,omitempty"`
	ConfigName        string      `json:"configName,omitempty"`
	ConfigContent     interface{} `json:"configContent,omitempty"`
	ConfigContentType string      `json:"configContentType,omitempty"`
}
