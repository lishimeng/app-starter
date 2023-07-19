package persistence

import (
	"github.com/beego/beego/v2/client/orm"
)

// DriverType RegisterModel
type DriverType orm.DriverType

type Driver struct {
	name string
	t    orm.DriverType
}

type BaseConfig struct {
	initDb     bool
	aliasName  string
	driver     Driver
	dataSource string
	params     []orm.DBOption
	models     []interface{}
}

func (b *BaseConfig) MaxIdle(n int) {
	if n > 0 {
		b.params = append(b.params, orm.MaxIdleConnections(n))
	}
}

func (b *BaseConfig) MaxConn(n int) {
	if n > 0 {
		b.params = append(b.params, orm.MaxOpenConnections(n))
	}
}

var DriverMysql = Driver{"mysql", orm.DRMySQL}
var DriverSqlite = Driver{"sqlite3", orm.DRSqlite}
var DriverOracle = Driver{"oracle", orm.DROracle}
var DriverPostgres = Driver{"postgres", orm.DRPostgres}
var DriverTiDB = Driver{"tidb", orm.DRTiDB}

var drivers map[string]orm.DriverType // driver缓存

func init() {
	drivers = make(map[string]orm.DriverType)
	err := autoRegisterDrivers()
	if err != nil {
		panic(err)
	}
}

func RegisterDriver(driver Driver) (err error) {

	name := driver.name
	d := driver.t
	if _, ok := drivers[name]; ok {
		return // 已经注册过了
	}
	err = orm.RegisterDriver(name, d)
	return
}

func autoRegisterDrivers() (err error) {
	var ds []Driver

	ds = append(ds, DriverPostgres) // 默认加载的driver

	for _, d := range ds {
		err = RegisterDriver(d)
		if err != nil {
			break
		}
	}
	return
}

func RegisterDataBase(init bool, aliasName, driverName, dataSource string, params ...orm.DBOption) (err error) {
	err = orm.RegisterDataBase(aliasName, driverName, dataSource, params...)
	if init {
		err = orm.RunSyncdb(aliasName, false, true)
	}
	return
}

func RegisterDatabase(config BaseConfig) (err error) {
	err = RegisterDataBase(config.initDb,
		config.aliasName,
		config.driver.name,
		config.dataSource,
		config.params...)
	return
}

func RegisterModels(models ...interface{}) {
	orm.RegisterModel(models...)
}

func (b *BaseConfig) RegisterModel(models ...interface{}) {
	b.models = append(b.models, models...)
}

func InitDatabase(config BaseConfig) (err error) {
	err = RegisterDataBase(config.initDb,
		config.aliasName,
		config.driver.name,
		config.dataSource,
		config.params...)
	return
}

func InitOrm(config BaseConfig) (err error) {

	// err = orm.RegisterDriver(config.driver.name, config.driver.t) // 不再初始化时加载driver， 改用默认的加载方式

	//if err != nil {
	//	return
	//}

	if len(config.models) > 0 {
		orm.RegisterModel(config.models...)
	}
	err = RegisterDataBase(config.initDb,
		config.aliasName,
		config.driver.name,
		config.dataSource,
		config.params...)
	return
}
