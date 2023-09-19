package cms

import (
	"github.com/lishimeng/app-starter"
	"time"
)

// WebSite 网站程序配置, 动态html输出
type WebSite struct {
	app.Pk
	Name      WebSiteName `orm:"column(name)"`      // 应用名称
	BaseUrl   string      `orm:"column(base_url)"`  // url前缀
	Copyright string      `orm:"column(copyright)"` // 版权
	Icp       string      `orm:"column(icp)"`       // icp
	Favicon   string      `orm:"column(favicon)"`   // favicon url
	Logo      string      `orm:"column(logo)"`      // logo名称
	Attr1     string      `orm:"column(attr1)"`     // 扩展
	Attr2     string      `orm:"column(attr2)"`     // 扩展
	Attr3     string      `orm:"column(attr3)"`     // 扩展
	app.TableChangeInfo
}

// SpaConfig SPA一般由vite控制, 使用localstorage管理数据
type SpaConfig struct {
	Id                int         `orm:"pk;auto;column(id)"`
	Name              WebSiteName `orm:"column(app_name);"`                           //应用名称
	ConfigName        string      `orm:"column(config_name);"`                        //配置字段名称
	ConfigContent     string      `orm:"column(config_content);default()"`            //配置字段内容
	ConfigContentType string      `orm:"column(config_content_Type);default(string)"` //配置字段内容类型
	CreateTime        time.Time   `orm:"auto_now_add;type(datetime);column(ctime)"`
}

// TableUnique 联合唯一约束
func (atc *SpaConfig) TableUnique() [][]string {
	return [][]string{
		{"Name", "ConfigName"},
	}
}

type ConfigContentType string

const (
	BooleanConfigContentType ConfigContentType = "boolean"
	NumberConfigContentType  ConfigContentType = "int"
	StringConfigContentType  ConfigContentType = "string"
)
