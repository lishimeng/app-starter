package token

import (
	"github.com/lishimeng/app-starter/cache"
)

// redisStorage 从redis获取token
type redisStorage struct {
	Conn        cache.C
	jwtProvider *JwtProvider
}

func NewRedisStorage(c cache.C, provider *JwtProvider) (s Storage) {
	s = &redisStorage{
		Conn:        c,
		jwtProvider: provider,
	}
	return
}

func (rs *redisStorage) Verify(key string) (p JwtPayload, err error) {
	ok := rs.Conn.Exists(key)
	if !ok {
		err = ErrInvalid
		return
	}
	vt, err := rs.jwtProvider.Verify([]byte(key))
	if err != nil {
		return
	}
	err = vt.Claims(&p)
	if err != nil {
		return
	}

	return
}
