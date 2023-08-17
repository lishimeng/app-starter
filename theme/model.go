package theme

import "time"

type AppWebView struct {
	Id         int       `orm:"pk;auto;column(id)"`
	AppName    string    `orm:"column(app_name);unique"` //应用名称
	ViewData   string    `orm:"column(view_data);null"`  //layout通用配置，json格式
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(ctime)"`
}

type AppThemeConfig struct {
	Id                int       `orm:"pk;auto;column(id)"`
	AppName           string    `orm:"column(app_name);"`                           //应用名称
	ConfigPage        string    `orm:"column(config_page);"`                        //配置所属页
	ConfigName        string    `orm:"column(config_name);"`                        //配置字段名称
	ConfigContent     string    `orm:"column(config_content);default()"`            //配置字段内容
	ConfigContentType string    `orm:"column(config_content_Type);default(string)"` //配置字段内容类型
	CreateTime        time.Time `orm:"auto_now_add;type(datetime);column(ctime)"`
}

// TableUnique 联合唯一约束
func (atc *AppThemeConfig) TableUnique() [][]string {
	return [][]string{
		{"AppName", "ConfigPage", "ConfigName"},
	}
}

type ConfigContentType string

const (
	BooleanConfigContentType ConfigContentType = "boolean"
	NumberConfigContentType  ConfigContentType = "int"
	StringConfigContentType  ConfigContentType = "string"
)
