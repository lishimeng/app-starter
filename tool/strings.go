package tool

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// 字符串工具包

// GetRandomString 随机字符串,hex表示.
//
// n:字节数,返回2n个hex字符
func GetRandomString(n int) (s string) {
	randBytes := make([]byte, n)
	_, _ = rand.Read(randBytes)
	s = fmt.Sprintf("%x", randBytes)
	return
}

func Join(delimiter string, s ...string) string {
	return strings.Join(s, delimiter)
}
