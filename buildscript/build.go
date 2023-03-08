package buildscript

import (
	"bytes"
	"fmt"
	"os"
	text "text/template"
)

// git update-index --chmod +x script.sh

type Param struct {
	Name  string
	Org   string
	Main  string
	HasUI bool
}

const (
	scriptName     = "build.sh"
	dockerFileName = "Dockerfile"
)

func Generate(name, org, mainPath string, hasUI bool) (err error) {

	fmt.Printf("hasUI:%+v\n", hasUI)
	p := Param{
		Name:  name,
		Org:   org,
		Main:  mainPath,
		HasUI: hasUI,
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
