package passwd

import (
	"crypto/sha512"
	"github.com/ZZMarquis/gm/sm3"
	"golang.org/x/crypto/bcrypt"
)

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
