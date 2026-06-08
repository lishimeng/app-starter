package persistence

import (
	"fmt"

	mysqldriver "gorm.io/driver/mysql"
	gormdb "gorm.io/gorm"
)

func init() {
	RegisterDialector(DriverMysql.Name, func(dsn string) gormdb.Dialector {
		return mysqldriver.Open(dsn)
	})
}

type MysqlConfig struct {
	InitDb    bool
	AliasName string
	UserName  string
	Password  string
	Host      string
	Port      int
	MaxIdle   int
	MaxConn   int
	DbName    string
}

func (c *MysqlConfig) Build() (b BaseConfig) {

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", c.UserName, c.Password, c.Host, c.Port, c.DbName)
	b = BaseConfig{
		DataSource: dataSource,
		AliasName:  c.AliasName,
		Driver:     DriverMysql,
		InitDb:     c.InitDb,
	}

	b.MaxIdle(c.MaxIdle)
	b.MaxConn(c.MaxConn)
	return
}
