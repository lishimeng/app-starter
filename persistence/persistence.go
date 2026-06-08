package persistence

import "fmt"

type Driver struct {
	Name string
}

var (
	DriverMysql    = Driver{"mysql"}
	DriverSqlite   = Driver{"sqlite3"}
	DriverOracle   = Driver{"oracle"}
	DriverPostgres = Driver{"postgres"}
	DriverTiDB     = Driver{"tidb"}
)

type BaseConfig struct {
	initDb     bool
	aliasName  string
	driver     Driver
	dataSource string
	maxIdle    int
	maxOpen    int
	models     []interface{}
}

func (b *BaseConfig) MaxIdle(n int) {
	if n > 0 {
		b.maxIdle = n
	}
}

func (b *BaseConfig) MaxConn(n int) {
	if n > 0 {
		b.maxOpen = n
	}
}

func (b *BaseConfig) RegisterModel(models ...interface{}) {
	b.models = append(b.models, models...)
}

func RegisterDataBase(init bool, aliasName, driverName, dataSource string, _ ...any) (err error) {
	cfg := BaseConfig{
		initDb:     init,
		aliasName:  aliasName,
		driver:     Driver{Name: driverName},
		dataSource: dataSource,
	}
	return RegisterDatabase(cfg)
}

func RegisterDatabase(config BaseConfig) (err error) {
	c := getConnector()
	if c == nil {
		return fmt.Errorf("persistence: no connector registered; import a backend such as persistence/beego")
	}

	alias := config.aliasName
	if alias == "" {
		alias = DefaultAlias
	}

	if len(config.models) > 0 {
		c.RegisterModels(config.models...)
	}

	opts := OpenOptions{
		Alias:    alias,
		Driver:   config.driver.Name,
		DSN:      config.dataSource,
		MaxIdle:  config.maxIdle,
		MaxOpen:  config.maxOpen,
		InitDB:   config.initDb,
	}

	session, err := c.Open(opts)
	if err != nil {
		return err
	}
	RegisterSession(alias, session)

	if config.initDb {
		err = c.Migrate(alias, config.models...)
	}
	return err
}

func RegisterModels(models ...interface{}) {
	c := getConnector()
	if c == nil {
		return
	}
	c.RegisterModels(models...)
}

func InitDatabase(config BaseConfig) (err error) {
	return RegisterDatabase(config)
}

func InitOrm(config BaseConfig) (err error) {
	return RegisterDatabase(config)
}
