package model

import "github.com/lishimeng/app-starter"

// BusinessConnector 业务连接器配置表 dp_business_connector
type BusinessConnector struct {
	app.Pk
	app.TableChangeInfo

	Code     string `orm:"column(code);unique"`
	Name     string `orm:"column(name)"`
	ConnType string `orm:"column(conn_type)"`
	Config   string `orm:"column(config)"`
	Enabled  int    `orm:"column(enabled);default(0)"`
	Desc     string `orm:"column(desc);null"`
}

func (BusinessConnector) TableName() string {
	return "dp_business_connector"
}
