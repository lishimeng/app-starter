package app

import "time"

type Pk struct {
	// ID
	Id int `orm:"pk;auto;column(id)"`
}

// Tenant
// 多租户
type Tenant struct {
	Org int `orm:"column(org)"` // org为tenant标记
}

// TenantPk
// 不可与 Pk 同时使用
type TenantPk struct {
	Pk
	Tenant
}

type OperatorInfo struct {
	CreateOperator int `orm:"column(coperator)"`
}

// OperatorChangeInfo
// 不可与 OperatorInfo 同时使用
type OperatorChangeInfo struct {
	OperatorInfo
	UpdateOperator int `orm:"column(moperator)"`
}

// TableChangeInfo
// 不可与 TableInfo 同时使用
type TableChangeInfo struct {
	// 状态
	Status int `orm:"column(status)"`
	// 创建时间
	TableInfo
	// 修改时间
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(mtime)"`
}

type TableInfo struct {
	// 创建时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(ctime)"`
}
