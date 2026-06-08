package persistencetest

import (
	"fmt"
	"testing"

	"github.com/lishimeng/app-starter/persistence"
)

func sqliteConfig(t *testing.T, initDb bool) persistence.BaseConfig {
	t.Helper()
	dsn := fmt.Sprintf("file:memdb_%s?mode=memory&cache=shared", t.Name())
	cfg := &persistence.SqliteConfig{
		Database:  dsn,
		AliasName: persistence.DefaultAlias,
		InitDb:    initDb,
	}
	return cfg.Build()
}

func registerTestModels() {
	persistence.RegisterModels(&TestRecord{})
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
