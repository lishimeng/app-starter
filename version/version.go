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

const (
	AppName = ""
	Version = ""
	Commit  = ""
	Build   = ""
)

func Print() {
	fmt.Println("***********************************")
	fmt.Println(AppName)
	fmt.Println(Version)
	fmt.Println(Commit)
	fmt.Println(Build)
	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOARCH)
	fmt.Println("***********************************")
}
