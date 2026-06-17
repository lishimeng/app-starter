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

// BaseConfig 数据库连接配置，由 PostgresConfig / MysqlConfig 等 Build 生成。
type BaseConfig struct {
	InitDb       bool
	SyncForce    bool
	SyncVerbose  bool
	AliasName    string
	Driver       Driver
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	Debug        bool
	Models       []any
	DriverOpts   any // 驱动专属选项，由对应 *Config.Build 填充
}

func (b *BaseConfig) MaxIdle(n int) {
	if n > 0 {
		b.MaxIdleConns = n
	}
}

func (b *BaseConfig) MaxConn(n int) {
	if n > 0 {
		b.MaxOpenConns = n
	}
}

func (b *BaseConfig) DebugLog(enable bool) {
	b.Debug = enable
}

func (b *BaseConfig) RegisterModel(models ...any) {
	b.Models = append(b.Models, models...)
}

func RegisterDataBase(init bool, aliasName, driverName, dataSource string, _ ...any) (err error) {
	cfg := BaseConfig{
		InitDb:     init,
		AliasName:  aliasName,
		Driver:     Driver{Name: driverName},
		DataSource: dataSource,
	}
	return RegisterDatabase(cfg)
}

func RegisterDatabase(config BaseConfig) (err error) {
	c := getConnector()
	if c == nil {
		return fmt.Errorf("persistence: no connector registered; call Install() and register dialectors")
	}

	if config.Driver.Name == DriverPostgres.Name {
		if err = validatePostgresDSN(config.DataSource); err != nil {
			return
		}
	}

	alias := config.AliasName
	if alias == "" {
		alias = DefaultAlias
	}

	if len(config.Models) > 0 {
		c.RegisterModels(config.Models...)
	}

	opts := OpenOptions{
		Alias:      alias,
		Driver:     config.Driver.Name,
		DSN:        config.DataSource,
		MaxIdle:    config.MaxIdleConns,
		MaxOpen:    config.MaxOpenConns,
		Debug:      config.Debug || isDebugEnabled(),
		InitDB:     config.InitDb,
		DriverOpts: config.DriverOpts,
	}

	session, err := c.Open(opts)
	if err != nil {
		return err
	}
	RegisterSession(alias, session)

	if config.InitDb {
		err = c.Migrate(alias, SyncOptions{
			Force:   config.SyncForce,
			Verbose: config.SyncVerbose,
		}, config.Models...)
	}
	return err
}

func RegisterModels(models ...any) {
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
