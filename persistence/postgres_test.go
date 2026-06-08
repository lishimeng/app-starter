package persistence

import "testing"

func TestInitPostgresOrm(t *testing.T) {

	alias := "default"
	c := PostgresConfig{
		AliasName: alias,
		InitDb:    false,
	}
	bc := c.Build()
	if bc.AliasName != alias {
		t.Fatal("alias name")
		return
	}
}

func TestValidatePostgresDSN(t *testing.T) {
	err := validatePostgresDSN("user= password= dbname= host= port=5432 sslmode=disable")
	if err == nil {
		t.Fatal("expected error for empty postgres config")
	}
	err = validatePostgresDSN("user=u password=p dbname=d host=h port=5432 sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
}
