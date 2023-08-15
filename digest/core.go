package digest

import (
	"crypto/sha512"
	"fmt"
	"github.com/ZZMarquis/gm/sm3"
	"github.com/lishimeng/app-starter/tool"
	"github.com/lishimeng/go-log"
	"golang.org/x/crypto/bcrypt"
)

func genPass(password string, nanoTime int64, digestFunc Hash, salting ...SaltingFunc) (r string) {
	plain := preEncode(password, nanoTime, salting...)

	if digestFunc == nil {
		digestFunc = defaultDigestFunc
	}

	bs, err := digestFunc([]byte(plain))
	if err != nil {
		return
	}
	r = string(bs)
	return
}

func verifyPass(encoded, password string, nanoTime int64, verifyFunc Verifier, salting ...SaltingFunc) (err error) {
	plain := preEncode(password, nanoTime, salting...)

	if verifyFunc == nil {
		verifyFunc = defaultVerifyFunc
	}

	err = verifyFunc([]byte(encoded), []byte(plain))
	if err != nil {
		return
	}
	return
}

func preEncode(password string, nanoTime int64, salting ...SaltingFunc) (r string) {
	var saltFunc SaltingFunc
	if len(salting) <= 0 {
		saltFunc = func(plaintext string) string {
			return fmt.Sprintf("%d.%s_%d", nanoTime, plaintext, nanoTime)
		}
	} else {
		saltFunc = salting[0]
	}

	r = saltFunc(password)
	return
}

func bcryptDigest(plain []byte) (bs []byte, err error) {
	bs, err = bcrypt.GenerateFromPassword(plain, bcrypt.DefaultCost)
	if err != nil {
		return
	}
	return
}
func bcryptVerify(encoded, plain []byte) (err error) {
	err = bcrypt.CompareHashAndPassword(encoded, plain)
	log.Info(err)
	if err != nil {
		err = ErrPasswordWrong
	}
	return
}

func sm3Digest(plain []byte) (bs []byte, err error) {
	p := sm3.New()
	_, err = p.Write(plain)
	if err != nil {
		return
	}
	s := hashHex(p.Sum(nil))
	bs = []byte(s)
	return
}
func sm3Verify(encoded, plain []byte) (err error) {
	p := sm3.New()
	_, err = p.Write(plain)
	if err != nil {
		return
	}
	s := hashHex(p.Sum(nil))
	if s != string(encoded) {
		err = ErrPasswordWrong
	}
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
func sha512Verify(encoded, plain []byte) (err error) {
	p := sha512.New()
	_, err = p.Write(plain)
	if err != nil {
		return
	}
	s := hashHex(p.Sum(nil))
	if s != string(encoded) {
		err = ErrPasswordWrong
	}
	return
}

func hashHex(bs []byte) string {
	return tool.BytesToHex(bs)
}
