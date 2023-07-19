package persistence

import (
	"fmt"
)

type PostgresConfig struct {
	InitDb    bool
	AliasName string
	UserName  string
	Password  string
	Host      string
	Port      int
	DbName    string
	MaxIdle   int
	MaxConn   int
	SSL       string
}

func (c *PostgresConfig) Build() (b BaseConfig) {

	ssl := "disable"
	if len(c.SSL) > 0 {
		ssl = c.SSL
	}
	dataSource := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", c.UserName, c.Password, c.DbName, c.Host, c.Port, ssl)
	b = BaseConfig{
		dataSource: dataSource,
		aliasName:  c.AliasName,
		driver:     DriverPostgres,
		initDb:     c.InitDb,
	}

	b.MaxIdle(c.MaxIdle)
	b.MaxConn(c.MaxConn)
	return
}
