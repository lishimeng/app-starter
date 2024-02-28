package token

import (
	"crypto/md5"
	"errors"
	"github.com/lishimeng/x/util"
)

var (
	ErrInvalid = errors.New("invalid token")
)

type ClientTye string

const (
	PC     ClientTye = "pc"
	App    ClientTye = "app"
	Pad    ClientTye = "pad"
	WeChat ClientTye = "wechat"
)

// 保存token值

type Payload struct {
	Key   string
	Value string
}

type Option struct {
}

// Storage 本地token存储。不负责创建token
type Storage interface {
	Verify(key string) (p JwtPayload, err error)
}

func FullVerify(key string) {

}

const (
	JwtTokenPrefix = "jwt_"
)

// Digest jwt hash摘要
func Digest(content []byte) (d string) {
	sh := md5.New()
	sh.Write(content)
	bs := sh.Sum(nil)
	d = JwtTokenPrefix + util.BytesToHex(bs)
	return
}
