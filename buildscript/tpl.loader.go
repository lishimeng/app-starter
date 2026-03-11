package buildscript

import (
	"embed"
	"io"
	"path"
)

//go:embed tpl
var tplFolder embed.FS

var baseDockerfileAlpineTpl string
var baseDockerfileUbuntuTpl string
var dockerfileTpl string
var buildShellV2Tpl string

func init() {
	initFile(&baseDockerfileAlpineTpl, "dockerfile.base.alpine.tpl")
	initFile(&baseDockerfileUbuntuTpl, "dockerfile.base.ubuntu.tpl")
	initFile(&dockerfileTpl, "dockerfile.tpl")
	initFile(&buildShellV2Tpl, "build.sh.v2.tpl")
}

func initFile(key *string, file string) {
	bs, err := readFile(file)
	if err != nil {
		panic(err)
	}
	*key = string(bs)
}

func readFile(name string) (bs []byte, err error) {
	uri := path.Join("tpl", name)
	file, err := tplFolder.Open(uri)
	if err != nil {
		return
	}
	bs, err = io.ReadAll(file)
	if err != nil {
		return
	}
	return
}
