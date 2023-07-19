package persistence

type SqliteConfig struct {
	Database  string
	AliasName string
	InitDb    bool
}

func (c *SqliteConfig) Build() (b BaseConfig) {

	b = BaseConfig{
		dataSource: c.Database,
		aliasName:  c.AliasName,
		driver:     DriverSqlite,
		initDb:     c.InitDb,
	}
	return
}
