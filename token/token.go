package token

import "errors"

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
