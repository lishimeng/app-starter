package repo

import "github.com/lishimeng/go-orm"

func Database(config persistence.BaseConfig, models ...interface{}) (err error) {

	// 此处配置默认数据库 alias=default
	//config.RegisterModel(models...)
	//err = persistence.InitOrm(config)

	persistence.RegisterModels(models...)
	err = persistence.RegisterDatabase(config)
	if err != nil {
		return
	}
	return
}
