package buildscript

import "testing"

func TestBuild(t *testing.T) {
	err := Generate("demo-app", "demo-org", "cmd/demo/main.go", false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenBaseDockerfile(t *testing.T) {
	err := GenerateBaseDockerfile()
	if err != nil {
		t.Fatal(err)
	}
}
