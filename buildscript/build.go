package buildscript

import (
	"bytes"
	"os"
	text "text/template"
)

// git update-index --chmod +x script.sh

type Param struct {
	Name string
	Org  string
	Main string
}

const (
	scriptName     = "build.sh"
	dockerFileName = "Dockerfile"
)

func Generate(name, org, mainPath string, hasUI bool) (err error) {

	p := Param{
		Name: name,
		Org:  org,
		Main: mainPath,
	}
	scriptContent, err := rendText(p, script)
	if err != nil {
		return
	}
	err = os.WriteFile(scriptName, []byte(scriptContent), 0644)

	var dockerContent string
	if hasUI {
		dockerContent, err = rendText(nil, dockerFileWithUI)
	} else {
		dockerContent, err = rendText(nil, dockerFile)
	}

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
