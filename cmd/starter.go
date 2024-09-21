package main

import (
	"flag"
	"fmt"
	"github.com/lishimeng/app-starter/buildscript"
	"os"
	"strings"
)

var (
	h        bool
	ns       string
	args     []string
	registry string
)

func init() {
	flag.Usage = usage
	flag.BoolVar(&h, "h", false, "help")
	flag.StringVar(&ns, "ns", "", "app namespace")
	flag.StringVar(&registry, "registry", "", "docker image registry")
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `App starter cli:
create a new app
cmd -org <org> <app_name:app_path:hasUI>
eg. cmd -org lishimeng name_a:cmd/subpath1:false name_b:cmd/subpath2:true
`)
	flag.PrintDefaults()
}

func main() {

	flag.Parse()
	args = flag.Args()

	if h {
		flag.Usage()
		return
	}
	if len(ns) == 0 || len(args) == 0 {
		flag.Usage()
		return
	}
	_main(args...)
}

func _main(components ...string) {
	var appInfo []buildscript.Application
	for _, content := range components {
		var name string
		var appPath string
		var hasUI = false
		kv := strings.Split(content, ":")
		if len(kv) > 1 {
			name = kv[0]
			appPath = kv[1]
		} else {
			fmt.Println("args error")
			flag.Usage()
			return
		}
		if len(kv) > 2 {
			tmp := kv[2]
			if tmp == "true" {
				hasUI = true
			}
		}
		info := buildscript.Application{Name: name, AppPath: appPath, HasUI: hasUI}
		appInfo = append(appInfo, info)
	}
	fmt.Println("Generate application")
	fmt.Println()
	fmt.Println("Namespace:", ns)
	if len(registry) > 0 {
		fmt.Println("Registry:", registry)
	}
	for _, info := range appInfo {
		fmt.Println("  app name:", info.Name)
		fmt.Println("  app path:", info.AppPath)
		fmt.Println("  hasUI:", info.HasUI)
		fmt.Println()
	}
	var projectInfo = buildscript.Project{
		Namespace: ns,
	}
	if len(registry) > 0 {
		projectInfo.ImageRegistry = registry
	}

	fmt.Println("Start generate...")
	err := buildscript.Generate(projectInfo, appInfo...)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Generate success")
	}
}
