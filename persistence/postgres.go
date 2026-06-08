package persistence

import (
	"fmt"
	"strings"

	pgdriver "gorm.io/driver/postgres"
	gormdb "gorm.io/gorm"
)

func init() {
	RegisterDialector(DriverPostgres.Name, func(dsn string) gormdb.Dialector {
		return pgdriver.Open(dsn)
	})
}

type PostgresConfig struct {
	InitDb         bool
	AliasName      string
	UserName       string
	Password       string
	Host           string
	Port           int
	DbName         string
	MaxIdle        int
	MaxConn        int
	SSL            string
	AdvancedConfig string // k=v a=b c=d m=n --> eg. sslcert/sslkey/sslrootcert/sslcrl
}

func (c *PostgresConfig) Build() (b BaseConfig) {

	ssl := "disable"
	if len(c.SSL) > 0 {
		ssl = c.SSL
	}
	dataSource := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", c.UserName, c.Password, c.DbName, c.Host, c.Port, ssl)
	if len(c.AdvancedConfig) > 0 {
		dataSource = fmt.Sprintf("%s %s", dataSource, c.AdvancedConfig)
	}
	b = BaseConfig{
		DataSource: dataSource,
		AliasName:  c.AliasName,
		Driver:     DriverPostgres,
		InitDb:     c.InitDb,
	}

	b.MaxIdle(c.MaxIdle)
	b.MaxConn(c.MaxConn)
	return
}

func validatePostgresDSN(dsn string) error {
	required := map[string]bool{"user": false, "dbname": false, "host": false}
	for _, part := range strings.Fields(dsn) {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		if _, exists := required[key]; exists && value != "" {
			required[key] = true
		}
	}

	var missing []string
	for key, ok := range required {
		if !ok {
			missing = append(missing, key)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf(
		"postgres config incomplete (missing %s): set DB_USER, DB_PASSWORD, DB_HOST, DB_DATABASE environment variables",
		strings.Join(missing, ", "),
	)
}
