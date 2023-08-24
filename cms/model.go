package cms

type WebSiteInfo struct {
	Name      WebSiteName // 应用名称
	BaseUrl   string      // url前缀
	Copyright string      // 版权
	Icp       string      // icp
	Favicon   string      // favicon url
	Logo      string      // logo名称
}

type WebSiteName string
