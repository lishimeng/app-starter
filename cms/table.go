package cms

import "github.com/lishimeng/app-starter"

type WebSite struct {
	app.Pk
	Name      WebSiteName `orm:"column(name)"`      // 应用名称
	BaseUrl   string      `orm:"column(base_url)"`  // url前缀
	Copyright string      `orm:"column(copyright)"` // 版权
	Icp       string      `orm:"column(icp)"`       // icp
	Favicon   string      `orm:"column(favicon)"`   // favicon url
	Logo      string      `orm:"column(logo)"`      // logo名称
	app.TableChangeInfo
}
