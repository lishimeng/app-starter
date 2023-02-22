package version

import (
	"fmt"
	"runtime"
)

/**

VERSION=$(git describe --tags $(git rev-list --tags --max-count=1))
COMMIT=$(git log --pretty=format:"%h" -1)
BUILD_TIME=$(date +%FT%T%z)
*/

var (
	AppName  = ""
	Version  = ""
	Commit   = ""
	Build    = ""
	Compiler = ""
)

func Print() {
	fmt.Println("***********************************")
	fmt.Printf("Name     :%s\n", AppName)
	fmt.Printf("Version  :%s\n", Version)
	fmt.Printf("Commit   :%s\n", Commit)
	fmt.Printf("Build    :%s\n", Build)
	fmt.Printf("Compiler :%s\n", Compiler)
	fmt.Printf("Runtime  :%s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("***********************************")
}
