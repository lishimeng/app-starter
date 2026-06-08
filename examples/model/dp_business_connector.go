package model

import "github.com/lishimeng/app-starter"

// BusinessConnector 业务连接器配置表 dp_business_connector
type BusinessConnector struct {
	app.Pk
	app.TableChangeInfo

	Code     string  `gorm:"column:code;uniqueIndex;not null;default:''"`
	Name     string  `gorm:"column:name;not null;default:''"`
	ConnType string  `gorm:"column:conn_type;not null;default:''"`
	Config   string  `gorm:"column:config;not null"`
	Enabled  int     `gorm:"column:enabled;not null;default:0"`
	Desc     *string `gorm:"column:desc"`
}

func (BusinessConnector) TableName() string {
	return "dp_business_connector"
}
