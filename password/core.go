package passwd

import (
	"fmt"
	"github.com/lishimeng/app-starter/tool"
)

func genPass(password string, nanoTime int64, digestFunc DigestFunc, salting ...SaltingFunc) (r string) {
	var saltFunc SaltingFunc
	if len(salting) <= 0 {
		saltFunc = func(plaintext string) string {
			return fmt.Sprintf("%d.%s_%d", nanoTime, plaintext, nanoTime)
		}
	} else {
		saltFunc = salting[0]
	}

	plain := saltFunc(password)

	if digestFunc == nil {
		digestFunc = defaultDigestFunc
	}

	bs, err := digestFunc([]byte(plain))
	if err != nil {
		return
	}
	r = tool.BytesToHex(bs)
	return
}
