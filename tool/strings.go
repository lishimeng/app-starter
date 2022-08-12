package tool

import (
	"fmt"
	"math/rand"
	"strings"
)

// 字符串工具包

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func Join(delimiter string, s ...string) string {
	return strings.Join(s, delimiter)
}
