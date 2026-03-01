package app

import (
	"os"
	"strings"
)

// InjectEnvStr key=value
func InjectEnvStr(s string) {
	list := strings.Split(s, "=")
	if len(list) != 2 {
		return
	}
	InjectEnv(list[0], list[1])
}

// InjectEnv key, value
func InjectEnv(key, value string) {
	_ = os.Setenv(key, value)
}
