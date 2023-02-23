package build

import (
	"bytes"
	"io/ioutil"
	text "text/template"
)

type Param struct {
	Name string
	Org  string
	Main string
}

const (
	scriptName     = "build.sh"
	dockerFileName = "Dockerfile"
)

func Generate(name, org, mainPath string) (err error) {

	p := Param{
		Name: name,
		Org:  org,
		Main: mainPath,
	}
	scriptContent, err := rendText(p, script)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(scriptName, []byte(scriptContent), 0644)

	dockerContent, err := rendText(nil, dockerFile)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(dockerFileName, []byte(dockerContent), 0644)
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
