package buildscript

import "testing"

func TestBuildHasUi(t *testing.T) {
	err := Generate(Project{
		ImageRegistry: "registry.my.domain.com",
		Namespace:     "proj111",
	}, Application{
		Name:    "Bala",
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
	err := Generate(Project{
		Namespace: "proj111",
	}, Application{
		Name:    "Bala",
		AppPath: "main_file_path",
		HasUI:   false,
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("ok")
	}
}

func TestGenBaseImage(t *testing.T) {
	err := GenerateBaseDockerfile()
	if err != nil {
		t.Fatal(err)
	}
}
