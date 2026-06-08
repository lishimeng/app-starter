package persistence

import (
	"github.com/glebarez/sqlite"
	gormdb "gorm.io/gorm"
)

func init() {
	RegisterDialector(DriverSqlite.Name, func(dsn string) gormdb.Dialector {
		return sqlite.Open(dsn)
	})
}

type SqliteConfig struct {
	Database  string
	AliasName string
	InitDb    bool
}

func (c *SqliteConfig) Build() (b BaseConfig) {

	b = BaseConfig{
		DataSource: c.Database,
		AliasName:  c.AliasName,
		Driver:     DriverSqlite,
		InitDb:     c.InitDb,
	}
	return
}
