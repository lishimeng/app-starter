package cms

import "github.com/lishimeng/app-starter"

func GetWebsiteFrame(name WebSiteName) (ws WebSiteInfo, err error) {

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
	return
}
