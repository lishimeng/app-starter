package token

import (
	"github.com/kataras/jwt"
	"github.com/lishimeng/app-starter/cache"
)

// redisStorage 从redis获取token
type redisStorage struct {
	Conn        cache.C
	jwtProvider *JwtProvider
}

func NewRedisStorage(c cache.C) (s Storage) {
	s = &redisStorage{Conn: c}
	return
}

func (rs *redisStorage) Verify(key string) (p JwtPayload, err error) {
	ut, err := jwt.Decode([]byte(key)) // 校验格式
	if err != nil {
		return
	}
	// TODO 校验时间

	tokenDigest := Digest([]byte(key)) // 计算摘要
	ok := rs.Conn.Exists(tokenDigest)  // 缓存检查
	if !ok {
		err = ErrInvalid
		return
	}

	err = ut.Claims(&p) // 取出payload
	if err != nil {
		return
	}

	return
}
