package persistence

import "testing"

func TestOpenWithoutDialector(t *testing.T) {
	if err := Install(); err != nil {
		t.Fatal(err)
	}
	_, err := defaultGormConnector.Open(OpenOptions{
		Alias:  "default",
		Driver: "postgres",
		DSN:    "unused",
	})
	if err == nil {
		t.Fatal("expected error when dialector is not registered")
	}
}

func TestInstallRegistersConnector(t *testing.T) {
	if err := Install(); err != nil {
		t.Fatal(err)
	}
}
