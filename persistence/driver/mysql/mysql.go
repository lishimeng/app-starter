// Package mysql provides MysqlConfig and registers the GORM mysql dialector on import.
package mysql

import (
	"fmt"

	"github.com/lishimeng/app-starter/persistence"
	mysqldriver "gorm.io/driver/mysql"
	gormdb "gorm.io/gorm"
)

func init() {
	persistence.RegisterDialector(persistence.DriverMysql.Name, openDialector)
}

func openDialector(opts persistence.OpenOptions) gormdb.Dialector {
	mo, _ := opts.DriverOpts.(*OpenOpts)
	if mo != nil && mo.hasAdvanced() {
		return mysqldriver.New(mysqldriver.Config{
			DSN:                       opts.DSN,
			DefaultStringSize:         mo.DefaultStringSize,
			DisableDatetimePrecision:  mo.DisableDatetimePrecision,
			DontSupportRenameIndex:    mo.DontSupportRenameIndex,
			DontSupportRenameColumn:   mo.DontSupportRenameColumn,
			SkipInitializeWithVersion: mo.SkipInitializeWithVersion,
		})
	}
	return mysqldriver.Open(opts.DSN)
}

// OpenOpts carries mysql driver options from Config.Build.
type OpenOpts struct {
	DefaultStringSize         uint
	DisableDatetimePrecision  bool
	DontSupportRenameIndex    bool
	DontSupportRenameColumn   bool
	SkipInitializeWithVersion bool
}

func (o *OpenOpts) hasAdvanced() bool {
	if o == nil {
		return false
	}
	return o.DefaultStringSize > 0 ||
		o.DisableDatetimePrecision ||
		o.DontSupportRenameIndex ||
		o.DontSupportRenameColumn ||
		o.SkipInitializeWithVersion
}

// Config mysql connection settings.
type Config struct {
	InitDb    bool
	AliasName string
	UserName  string
	Password  string
	Host      string
	Port      int
	MaxIdle   int
	MaxConn   int
	DbName    string
	Charset   string
	DisableParseTime bool
	Loc       string
	DefaultStringSize         uint
	DisableDatetimePrecision  bool
	DontSupportRenameIndex    bool
	DontSupportRenameColumn   bool
	SkipInitializeWithVersion bool
}

func (c *Config) Build() (b persistence.BaseConfig) {
	charset := c.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	parseTime := "True"
	if c.DisableParseTime {
		parseTime = "False"
	}
	loc := c.Loc
	if loc == "" {
		loc = "Local"
	}
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		c.UserName, c.Password, c.Host, c.Port, c.DbName, charset, parseTime, loc)

	b = persistence.BaseConfig{
		DataSource: dataSource,
		AliasName:  c.AliasName,
		Driver:     persistence.DriverMysql,
		InitDb:     c.InitDb,
		DriverOpts: &OpenOpts{
			DefaultStringSize:         c.DefaultStringSize,
			DisableDatetimePrecision:  c.DisableDatetimePrecision,
			DontSupportRenameIndex:    c.DontSupportRenameIndex,
			DontSupportRenameColumn:   c.DontSupportRenameColumn,
			SkipInitializeWithVersion: c.SkipInitializeWithVersion,
		},
	}
	b.MaxIdle(c.MaxIdle)
	b.MaxConn(c.MaxConn)
	return
}
