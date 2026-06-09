package setup

import (
	"os"

	"github.com/lishimeng/app-starter/persistence/driver/postgres"
)

func PostgresConfig() *postgres.Config {
	cfg := postgres.Config{
		UserName:       os.Getenv("DB_USER"),
		Password:       os.Getenv("DB_PASSWORD"),
		Host:           os.Getenv("DB_HOST"),
		Port:           5432,
		DbName:         os.Getenv("DB_DATABASE"),
		InitDb:         false,
		AliasName:      "default",
		TimeZone:       "Asia/Shanghai",
	}
	return &cfg
}
