package digest

import (
	"crypto/sha512"
	"fmt"
	"github.com/ZZMarquis/gm/sm3"
	"github.com/lishimeng/app-starter/tool"
	"golang.org/x/crypto/bcrypt"
)

func genPass(password string, nanoTime int64, digestFunc Hash, salting ...SaltingFunc) (r string) {
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

func bcryptDigest(plain []byte) (bs []byte, err error) {
	bs, err = bcrypt.GenerateFromPassword(plain, bcrypt.DefaultCost)
	if err != nil {
		return
	}
	return
}
func sm3Digest(plain []byte) (bs []byte, err error) {
	p := sm3.New()
	_, err = p.Write(plain)
	if err != nil {
		return
	}
	bs = p.Sum(nil)
	return
}

func sha512Digest(plain []byte) (bs []byte, err error) {
	p := sha512.New()
	_, err = p.Write(plain)
	if err != nil {
		return
	}
	bs = p.Sum(nil)
	return
}
