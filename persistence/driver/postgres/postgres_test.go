package postgres

import (
	"strings"
	"testing"

	"github.com/lishimeng/app-starter/persistence"
)

func TestConfigBuild(t *testing.T) {
	bc := (&Config{
		AliasName:            "default",
		UserName:             "u",
		Password:             "p",
		Host:                 "h",
		Port:                 5432,
		DbName:               "d",
		TimeZone:             "Asia/Shanghai",
		PreferSimpleProtocol: true,
	}).Build()
	if bc.AliasName != "default" {
		t.Fatal("alias")
	}
	if !strings.Contains(bc.DataSource, "TimeZone=Asia/Shanghai") {
		t.Fatalf("dsn: %s", bc.DataSource)
	}
	opts, ok := bc.DriverOpts.(*OpenOpts)
	if !ok || !opts.PreferSimpleProtocol {
		t.Fatal("driver opts")
	}
	if openDialector(persistence.OpenOptions{DSN: bc.DataSource, DriverOpts: bc.DriverOpts}) == nil {
		t.Fatal("dialector")
	}
}
