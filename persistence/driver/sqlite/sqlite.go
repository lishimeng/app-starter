// Package sqlite provides SqliteConfig and registers the GORM sqlite dialector on import.
package sqlite

import (
	"github.com/glebarez/sqlite"
	"github.com/lishimeng/app-starter/persistence"
	gormdb "gorm.io/gorm"
)

func init() {
	persistence.RegisterDialector(persistence.DriverSqlite.Name, func(opts persistence.OpenOptions) gormdb.Dialector {
		return sqlite.Open(opts.DSN)
	})
}

// Config sqlite connection settings.
type Config struct {
	Database  string
	AliasName string
	InitDb    bool
}

func (c *Config) Build() (b persistence.BaseConfig) {
	b = persistence.BaseConfig{
		DataSource: c.Database,
		AliasName:  c.AliasName,
		Driver:     persistence.DriverSqlite,
		InitDb:     c.InitDb,
	}
	return
}
