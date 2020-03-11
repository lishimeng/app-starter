package repo

import "github.com/lishimeng/go-libs/persistence"

func Database(config persistence.BaseConfig, models ...interface{}) (ormContext *persistence.OrmContext, err error) {

	config.RegisterModel(models...)
	orm, err := persistence.InitOrm(config)
	if err != nil {
		return
	}
	ormContext = &orm
	return
}
