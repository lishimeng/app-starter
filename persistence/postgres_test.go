package persistence

import "testing"

func TestInitPostgresOrm(t *testing.T) {

	alias := "default"
	c := PostgresConfig{
		AliasName: alias,
		InitDb:    false,
	}
	bc := c.Build()
	if bc.aliasName != alias {
		t.Fatal("alias name")
		return
	}
}
