package beego

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/lishimeng/app-starter/persistence"
)

func init() {
	_ = RegisterDrivers(DriverPostgres)
	persistence.SetConnector(&connector{})
	persistence.SetConditionFactory(func() persistence.Condition {
		return newCondition()
	})
	persistence.SetFallbackSessionFactory(func(alias string) persistence.Session {
		if alias == "" {
			alias = persistence.DefaultAlias
		}
		if alias == persistence.DefaultAlias {
			return newSession(alias, orm.NewOrm())
		}
		return newSession(alias, orm.NewOrmUsingDB(alias))
	})
}
