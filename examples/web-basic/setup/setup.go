package setup

import (
	"os"

	"github.com/lishimeng/app-starter/persistence/driver/postgres"
)

func WebPort() string {
	env := os.Getenv("WEB_PORT")
	if len(env) > 0 {
		return env
	}
	return ":9527"
}

func PostgresConfig() *postgres.Config {
	cfg := postgres.Config{
		UserName:    os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		Host:        os.Getenv("DB_HOST"),
		Port:        5432,
		DbName:      os.Getenv("DB_DATABASE"),
		InitDb:      true,
		AliasName:   "default",
		TimeZone:    "Asia/Shanghai",
		SyncVerbose: true,
	}
	return &cfg
}
