package passwd

import (
	"fmt"
	"github.com/lishimeng/app-starter/tool"
)

func Generate(plaintext string, nanoTime int64) (r string) {
	r = genPass(plaintext, nanoTime)
	return
}

func Verify(plaintext string, encodedPassword string, nanoTime int64) (r bool) {
	encoded := genPass(plaintext, nanoTime)
	r = encoded == encodedPassword
	return
}

func genPass(password string, nanoTime int64) (r string) {
	s := nanoTime
	plain := fmt.Sprintf("%d.%s_%d", s, password, s)
	bs, err := gen([]byte(plain))
	if err != nil {
		return
	}
	r = tool.BytesToHex(bs)
	return
}
