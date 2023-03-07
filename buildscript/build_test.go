package buildscript

import "testing"

func TestBuild(t *testing.T) {
	err := Generate("demo-app", "demo-org", "cmd/demo/main.go", false)
	if err != nil {
		t.Fatal(err)
	}
}
