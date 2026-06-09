package mysql

import (
	"strings"
	"testing"
)

func TestConfigBuildDefaults(t *testing.T) {
	bc := (&Config{
		UserName: "u",
		Password: "p",
		Host:     "h",
		Port:     3306,
		DbName:   "d",
	}).Build()
	if !strings.Contains(bc.DataSource, "charset=utf8mb4") ||
		!strings.Contains(bc.DataSource, "parseTime=True") ||
		!strings.Contains(bc.DataSource, "loc=Local") {
		t.Fatalf("unexpected dsn: %s", bc.DataSource)
	}
}
