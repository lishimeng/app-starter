// Package postgres provides PostgresConfig and registers the GORM postgres dialector on import.
package postgres

import (
	"fmt"
	"strings"

	"github.com/lishimeng/app-starter/persistence"
	pgdriver "gorm.io/driver/postgres"
	gormdb "gorm.io/gorm"
)

func init() {
	persistence.RegisterDialector(persistence.DriverPostgres.Name, openDialector)
}

func openDialector(opts persistence.OpenOptions) gormdb.Dialector {
	po, _ := opts.DriverOpts.(*OpenOpts)
	if po != nil && po.PreferSimpleProtocol {
		return pgdriver.New(pgdriver.Config{
			DSN:                  opts.DSN,
			PreferSimpleProtocol: true,
		})
	}
	return pgdriver.Open(opts.DSN)
}

// OpenOpts carries postgres driver options from Config.Build.
type OpenOpts struct {
	PreferSimpleProtocol bool
}

// Config postgres connection settings.
type Config struct {
	InitDb               bool
	SyncForce            bool
	SyncVerbose          bool
	AliasName            string
	UserName       string
	Password       string
	Host           string
	Port           int
	DbName         string
	MaxIdle        int
	MaxConn        int
	SSL            string
	TimeZone       string
	AdvancedConfig string
	PreferSimpleProtocol bool
}

func (c *Config) Build() (b persistence.BaseConfig) {
	ssl := "disable"
	if len(c.SSL) > 0 {
		ssl = c.SSL
	}
	dataSource := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", c.UserName, c.Password, c.DbName, c.Host, c.Port, ssl)
	if c.TimeZone != "" &&
		!strings.Contains(dataSource, "TimeZone=") &&
		!strings.Contains(strings.ToLower(dataSource), "timezone=") {
		dataSource = fmt.Sprintf("%s TimeZone=%s", dataSource, c.TimeZone)
	}
	if len(c.AdvancedConfig) > 0 {
		dataSource = fmt.Sprintf("%s %s", dataSource, c.AdvancedConfig)
	}
	b = persistence.BaseConfig{
		DataSource:  dataSource,
		AliasName:   c.AliasName,
		Driver:      persistence.DriverPostgres,
		InitDb:      c.InitDb,
		SyncForce:   c.SyncForce,
		SyncVerbose: c.SyncVerbose,
		DriverOpts: &OpenOpts{
			PreferSimpleProtocol: c.PreferSimpleProtocol,
		},
	}
	b.MaxIdle(c.MaxIdle)
	b.MaxConn(c.MaxConn)
	return
}
