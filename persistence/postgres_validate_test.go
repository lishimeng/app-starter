package persistence

import "testing"

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
