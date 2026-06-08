package app

import "time"

type Pk struct {
	// ID
	Id int `gorm:"primaryKey;autoIncrement;column:id"`
}

// Pk64 主键为 bigint / bigserial 时使用；不可与 Pk 同时使用
type Pk64 struct {
	Id int64 `gorm:"primaryKey;autoIncrement;column:id"`
}

// Tenant 多租户
type Tenant struct {
	Org int `gorm:"column:org"` // org为tenant标记
}

// TenantPk 不可与 Pk / Pk64 同时使用
type TenantPk struct {
	Pk
	Tenant
}

// TenantPk64 不可与 Pk / Pk64 同时使用
type TenantPk64 struct {
	Pk64
	Tenant
}

type OperatorInfo struct {
	CreateOperator int `gorm:"column:coperator"`
}

// OperatorChangeInfo 不可与 OperatorInfo 同时使用
type OperatorChangeInfo struct {
	OperatorInfo
	UpdateOperator int `gorm:"column:moperator"`
}

// TableChangeInfo 不可与 TableInfo 同时使用
type TableChangeInfo struct {
	Status int `gorm:"column:status;default:0"`
	TableInfo
	UpdateTime time.Time `gorm:"autoUpdateTime;column:mtime"`
}

type TableInfo struct {
	CreateTime time.Time `gorm:"autoCreateTime;column:ctime"`
}

const (
	Disable = iota
	Enable
)
