package buildscript

import "testing"

func TestBuild(t *testing.T) {
	err := Generate("demo-org", true, App("a", "cmd/a"), App("b", "cmd/b"))
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
