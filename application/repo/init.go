package repo

import "github.com/lishimeng/app-starter/persistence"

func Database(config persistence.BaseConfig, views []any, models ...any) (err error) {
	if err = persistence.Install(); err != nil {
		return
	}
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
