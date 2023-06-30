package buildscript

import (
	"bytes"
	"os"
	text "text/template"
)

// git update-index --chmod +x script.sh

type Param struct {
	Org          string
	HasUI        bool
	Applications []Application
}
type Application struct {
	Name    string
	AppPath string
}

const (
	scriptName     = "build.sh"
	dockerFileName = "Dockerfile"
)

func App(name, appPath string) Application {
	return Application{
		Name:    name,
		AppPath: appPath,
	}
}

func Generate(org string, hasUI bool, apps ...Application) (err error) {

	p := Param{
		Org:          org,
		HasUI:        hasUI,
		Applications: apps,
	}

	scriptContent, err := rendText(p, script)
	if err != nil {
		return
	}
	err = os.WriteFile(scriptName, []byte(scriptContent), 0644)

	dockerContent, err := rendText(p, dockerFile)

	if err != nil {
		return
	}

	err = os.WriteFile(dockerFileName, []byte(dockerContent), 0644)
	return
}

func rendText(data interface{}, temp string) (content string, err error) {
	t, err := text.New("_").Parse(temp)
	if err != nil {
		return
	}
	w := new(bytes.Buffer)
	err = t.Execute(w, data)
	if err != nil {
		return
	}
	content = w.String()
	return
}

func GenerateBaseDockerfile() (err error) {

	dockerContent, err := rendText(nil, baseDockerFile)

	if err != nil {
		return
	}

	err = os.WriteFile(dockerFileName, []byte(dockerContent), 0644)
	return
}
