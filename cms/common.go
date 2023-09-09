package cms

const (
	pageThemeCacheKeyTpl = "theme_app-%s"
)

type SpaResp struct {
	Code    interface{} `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type SpaConfigInfo struct {
	Id                int         `json:"id,omitempty"`
	Name              string      `json:"name,omitempty"`
	ConfigPage        string      `json:"configPage,omitempty"`
	ConfigName        string      `json:"configName,omitempty"`
	ConfigContent     interface{} `json:"configContent,omitempty"`
	ConfigContentType string      `json:"configContentType,omitempty"`
}
