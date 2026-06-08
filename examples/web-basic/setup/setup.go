package setup

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/lishimeng/app-starter/persistence"
)

var ()

func PostgresConfig() *persistence.PostgresConfig {
	cfg := persistence.PostgresConfig{
		UserName:       os.Getenv("DB_USER"),
		Password:       os.Getenv("DB_PASSWORD"),
		Host:           os.Getenv("DB_HOST"),
		Port:           5432,
		DbName:         os.Getenv("DB_DATABASE"),
		InitDb:         false,
		AliasName:      "default",
		AdvancedConfig: "timezone=Asia/Shanghai",
	}
	return &cfg
}
