package persistence

import (
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
)

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
		dataSource: dataSource,
		aliasName:  c.AliasName,
		driver:     DriverMysql,
		initDb:     c.InitDb,
	}

	b.MaxIdle(c.MaxIdle)
	b.MaxConn(c.MaxConn)
	return
}
