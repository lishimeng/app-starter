package repo

import "github.com/lishimeng/app-starter/persistence"

func Database(config persistence.BaseConfig, views []any, models ...any) (err error) {

	// 此处配置默认数据库 alias=default
	//config.RegisterModel(models...)
	//err = persistence.InitOrm(config)

	if len(models) > 0 {
		persistence.RegisterModels(models...)
	}
	err = persistence.RegisterDatabase(config)
	if err != nil {
		return
	}
	if len(views) > 0 {
		persistence.RegisterModels(views...)
	}
	return
}
