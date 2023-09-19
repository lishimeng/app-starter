package cms

import (
	"github.com/lishimeng/app-starter"
)

var c WebSiteInfo

func fromConfig(info WebSiteInfo) {
	c = info
}

func getWebsiteFromDb(name WebSiteName) (ws WebSiteInfo, err error) {

	var website WebSite
	err = app.GetOrm().Context.
		QueryTable(new(WebSite)).
		Filter("Name", name).
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

	ws.Attr1 = website.Attr1
	ws.Attr2 = website.Attr2
	ws.Attr3 = website.Attr3
	return
}

func getDefaultConfig() (ws WebSiteInfo) {

	ws.Name = "app"
	ws.BaseUrl = "http://localhost"
	ws.Copyright = ""
	ws.Icp = ""
	ws.Favicon = ""
	ws.Logo = ""
	return
}
