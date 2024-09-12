package buildscript

import "testing"

func TestBuildHasUi(t *testing.T) {
	err := Generate("ORG_DEMO", Application{
		Name:    "AppName",
		AppPath: "main_file_path",
		HasUI:   true,
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("ok")
	}
}

func TestBuildNoUi(t *testing.T) {
	err := Generate("ORG_DEMO", Application{
		Name:    "AppName",
		AppPath: "main_file_path",
		HasUI:   false,
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("ok")
	}
}
