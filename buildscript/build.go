package buildscript

import (
	"bytes"
	"fmt"
	"os"
	text "text/template"

	"github.com/lishimeng/go-log"
)

// git update-index --chmod +x script.sh

// runtime image config
var (
	NodeImageVersion    = "node:24"
	GolangImageVersion  = "golang:1.24.10"
	RuntimeImageVersion = "lishimeng/alpine:3.17"
	ImageSubNamespace   = "library"
)

type Project struct {
	ImageRegistry  string // 镜像库,默认empty
	Namespace      string // namespace
	NamespaceAlias string //  namespace别名,不为空时替换namespace
}

func (p *Project) Build() *Project {
	if len(p.NamespaceAlias) > 0 {
		p.Namespace = p.NamespaceAlias
	}
	return p
}

type Param struct {
	Pro          Project
	Applications []Application
}
type DockerParam struct {
	Pro               Project
	App               Application
	BuildImageVersion ImageVersion
}
type Application struct {
	Name      string
	NameAlias string // name别名, 不为空时替换name
	AppPath   string
	HasUI     bool
}

func (a *Application) Build() *Application {
	if len(a.NameAlias) > 0 {
		a.Name = a.NameAlias
	}
	return a
}

type ImageVersion struct {
	Node         string // UI编译镜像
	Golang       string // golang编译镜像
	Runtime      string // 运行时镜像
	SubNamespace string // 镜像存放的二级地址,registry设置后生效 不空:{registry}/{subNamespace} 空:{registry}
}

const (
	localBuildScriptName = "build_local.sh"
	scriptName           = "build.sh"
	dockerFileName       = "Dockerfile"
)

func Generate(p Project, apps ...Application) (err error) {

	err = createShell(p, apps...)
	if err != nil {
		return
	}

	//err = createLocalBuildShell(p, apps...)
	//if err != nil {
	//	return
	//}

	err = createDockers(p, apps...)
	return
}

func createLocalBuildShell(pro Project, apps ...Application) (err error) {
	var compiledApp []Application
	for _, app := range apps {
		compiledApp = append(compiledApp, *app.Build())
	}
	var compiledPro = *pro.Build()
	p := Param{
		Pro:          compiledPro,
		Applications: compiledApp,
	}
	if len(p.Pro.ImageRegistry) > 0 {
		p.Pro.Namespace = fmt.Sprintf("%s/%s", p.Pro.ImageRegistry, p.Pro.Namespace)
	}
	scriptContent, err := rendText(p, localBuildScript)
	if err != nil {
		return
	}
	err = os.WriteFile(localBuildScriptName, []byte(scriptContent), 0644)
	if err != nil {
		log.Info(err)
		return
	}
	return
}

func createShell(pro Project, apps ...Application) (err error) {
	var compiledApp []Application
	for _, app := range apps {
		compiledApp = append(compiledApp, *app.Build())
	}
	var compiledPro = *pro.Build()
	p := Param{
		Pro:          compiledPro,
		Applications: compiledApp,
	}
	if len(p.Pro.ImageRegistry) > 0 {
		p.Pro.Namespace = fmt.Sprintf("%s/%s", p.Pro.ImageRegistry, p.Pro.Namespace)
	}
	scriptContent, err := rendText(p, scriptV2)
	if err != nil {
		return
	}
	err = os.WriteFile(scriptName, []byte(scriptContent), 0644)
	if err != nil {
		log.Info(err)
		return
	}
	return
}

func createDockers(p Project, apps ...Application) (err error) {
	p = *p.Build()
	for _, app := range apps {
		app = *app.Build()
		err = createDocker(p, app, ImageVersion{
			Node:         NodeImageVersion,
			Golang:       GolangImageVersion,
			Runtime:      RuntimeImageVersion,
			SubNamespace: ImageSubNamespace,
		})
		if err != nil {
			return
		}
	}
	return
}

func createDocker(pro Project, app Application, version ImageVersion) (err error) {
	p := DockerParam{
		Pro:               pro,
		App:               app,
		BuildImageVersion: version,
	}
	if len(pro.ImageRegistry) > 0 {
		var registry = pro.ImageRegistry
		if len(version.SubNamespace) > 0 {
			registry = fmt.Sprintf("%s/%s", registry, version.SubNamespace)
		}
		p.BuildImageVersion.Node = fmt.Sprintf("%s/%s", registry, p.BuildImageVersion.Node)
		p.BuildImageVersion.Golang = fmt.Sprintf("%s/%s", registry, p.BuildImageVersion.Golang)
		p.BuildImageVersion.Runtime = fmt.Sprintf("%s/%s", registry, p.BuildImageVersion.Runtime)
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

func GenerateBaseDockerfile(category ...string) (err error) {

	tpl := baseDockerfileAlpine
	if len(category) > 0 && category[0] == "ubuntu" {
		tpl = baseDockerfileUbuntu
	}
	dockerContent, err := rendText(nil, tpl)

	if err != nil {
		return
	}

	err = os.WriteFile(dockerFileName, []byte(dockerContent), 0644)
	return
}
