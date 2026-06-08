package beego

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/lishimeng/app-starter/persistence"
)

type connector struct{}

func (c *connector) Open(opts persistence.OpenOptions) (persistence.Session, error) {
	alias := opts.Alias
	if alias == "" {
		alias = persistence.DefaultAlias
	}

	params := make([]orm.DBOption, 0, 2)
	if opts.MaxIdle > 0 {
		params = append(params, orm.MaxIdleConnections(opts.MaxIdle))
	}
	if opts.MaxOpen > 0 {
		params = append(params, orm.MaxOpenConnections(opts.MaxOpen))
	}
	for _, p := range opts.DBParams {
		if opt, ok := p.(orm.DBOption); ok {
			params = append(params, opt)
		}
	}

	err := orm.RegisterDataBase(alias, opts.Driver, opts.DSN, params...)
	if err != nil {
		return nil, err
	}

	var o orm.Ormer
	if alias == persistence.DefaultAlias {
		o = orm.NewOrm()
	} else {
		o = orm.NewOrmUsingDB(alias)
	}

	if opts.Debug {
		orm.Debug = true
	}

	return newSession(alias, o), nil
}

func (c *connector) Migrate(alias string, models ...any) error {
	if alias == "" {
		alias = persistence.DefaultAlias
	}
	if len(models) > 0 {
		orm.RegisterModel(models...)
	}
	return orm.RunSyncdb(alias, false, true)
}

func (c *connector) RegisterModels(models ...any) {
	if len(models) == 0 {
		return
	}
	orm.RegisterModel(models...)
}

// RegisterDriver registers a beego ORM driver by name.
func RegisterDriver(name string, t orm.DriverType) error {
	return orm.RegisterDriver(name, t)
}

// Driver bundles a driver name and beego driver type.
type Driver struct {
	Name string
	Type orm.DriverType
}

var (
	DriverMysql    = Driver{"mysql", orm.DRMySQL}
	DriverSqlite   = Driver{"sqlite3", orm.DRSqlite}
	DriverOracle   = Driver{"oracle", orm.DROracle}
	DriverPostgres = Driver{"postgres", orm.DRPostgres}
	DriverTiDB     = Driver{"tidb", orm.DRTiDB}
)

var registeredDrivers = make(map[string]struct{})

// RegisterDrivers registers the given drivers once.
func RegisterDrivers(drivers ...Driver) error {
	for _, d := range drivers {
		if _, ok := registeredDrivers[d.Name]; ok {
			continue
		}
		if err := RegisterDriver(d.Name, d.Type); err != nil {
			return fmt.Errorf("register driver %s: %w", d.Name, err)
		}
		registeredDrivers[d.Name] = struct{}{}
	}
	return nil
}
