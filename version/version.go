package version

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
