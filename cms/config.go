package cms

import (
	"github.com/lishimeng/app-starter"
)

var c WebSiteInfo

func fromConfig(info WebSiteInfo) {
	info = c
}

func getWebsiteFromDb(name WebSiteName) (ws WebSiteInfo, err error) {

	var website WebSite
	err = app.GetOrm().Context.
		QueryTable(new(WebSite)).
		Filter("AppName", name).
		Filter("Status", app.Enable).
		One(&website)
	if err != nil {
		return
	}
	ws.Name = name
	ws.BaseUrl = website.BaseUrl
	ws.Copyright = website.Copyright
	ws.Icp = website.Icp
	ws.Favicon = website.Favicon
	ws.Logo = website.Logo
	return
}

func getDefaultConfig() (ws WebSiteInfo) {

	ws.Name = ""
	ws.BaseUrl = ""
	ws.Copyright = ""
	ws.Icp = ""
	ws.Favicon = ""
	ws.Logo = ""
	return
}
