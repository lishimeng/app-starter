package buildscript

import (
	"bytes"
	"os"
	text "text/template"
)

// git update-index --chmod +x script.sh

type Param struct {
	Org          string
	Applications []Application
}
type DockerParam struct {
	Org string
	App Application
}
type Application struct {
	Name    string
	AppPath string
	HasUI   bool
}

const (
	scriptName     = "build.sh"
	dockerFileName = "Dockerfile"
)

func Generate(org string, apps ...Application) (err error) {

	err = createShell(org, apps...)
	if err != nil {
		return
	}

	err = createDockers(org, apps...)
	return
}

func createShell(org string, apps ...Application) (err error) {
	p := Param{
		Org:          org,
		Applications: apps,
	}
	scriptContent, err := rendText(p, script)
	if err != nil {
		return
	}
	err = os.WriteFile(scriptName, []byte(scriptContent), 0644)
	return
}

func createDockers(org string, apps ...Application) (err error) {
	for _, app := range apps {
		err = createDocker(org, app)
		if err != nil {
			return
		}
	}
	return
}

func createDocker(org string, app Application) (err error) {
	p := DockerParam{
		Org: org,
		App: app,
	}
	dockerContent, err := rendText(p, dockerFile)
	if err != nil {
		return
	}

	dir := app.AppPath
	err = os.MkdirAll(dir, 755)
	if err != nil {
		return
	}
	err = os.WriteFile(dir+"/"+dockerFileName, []byte(dockerContent), 0644)
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
