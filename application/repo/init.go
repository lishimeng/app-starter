package repo

import "github.com/lishimeng/go-orm"

func Database(config persistence.BaseConfig, models ...interface{}) (err error) {

	config.RegisterModel(models...)
	err = persistence.InitOrm(config)
	if err != nil {
		return
	}
	return
}
