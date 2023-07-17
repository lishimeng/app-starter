package token

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kataras/jwt"
	"time"
)

type JwtPayload struct {
	Org    string `json:"org,omitempty"`
	Dept   string `json:"dept,omitempty"`
	Uid    string `json:"uid,omitempty"`
	Client string `json:"client,omitempty"`
	Scope  string `json:"scope,omitempty"` // 逗号分隔a,b,c,d,.....
}

func (jp *JwtPayload) WithClient(clientType ClientTye) *JwtPayload {
	switch clientType {
	case App:
	case Pad:
	case PC:
	case WeChat:
	default:
		return jp
	}
	jp.Client = string(clientType)
	return jp
}

type JwtProvider struct {
	alg         jwt.Alg
	signKey     jwt.PrivateKey
	verifyKey   jwt.PublicKey
	issuer      string
	defaultTTL  time.Duration
	checkIssuer bool
}

type JwtBuildOption func(provider *JwtProvider)

// HS256：bit 长度要>=256，即字节长度>=32
// HS384：bit 长度要>=384，即字节长度>=48
// HS512：bit 长度要>=512，即字节长度>=64
var jwtAllAlg = []jwt.Alg{
	jwt.NONE,
	jwt.HS256,
	jwt.HS384,
	jwt.HS512,
	jwt.RS256,
	jwt.RS384,
	jwt.RS512,
	jwt.PS256,
	jwt.PS384,
	jwt.PS512,
	jwt.ES256,
	jwt.ES384,
	jwt.ES512,
	jwt.EdDSA,
}

var WithAlg = func(alg string) JwtBuildOption {
	return func(provider *JwtProvider) {
		for _, a := range jwtAllAlg {
			if a.Name() == alg {
				provider.alg = a
				break
			}
		}
	}
}

// WithIssuer 设置issuer, 如不设置,将不能创建jwt, 校验的时候不会检查issuer是否相同
var WithIssuer = func(issuer string) JwtBuildOption {
	return func(provider *JwtProvider) {
		provider.issuer = issuer
		provider.checkIssuer = true
	}
}

// WithKey 设置秘钥
var WithKey = func(signKey interface{}, verifyKey interface{}) JwtBuildOption {
	return func(provider *JwtProvider) {
		provider.signKey = signKey
		provider.verifyKey = verifyKey
	}
}

// WithDefaultTTL 设置ttl
var WithDefaultTTL = func(ttl time.Duration) JwtBuildOption {
	return func(provider *JwtProvider) {
		provider.defaultTTL = ttl
	}
}

func NewJwtProvider(options ...JwtBuildOption) (jp *JwtProvider) {
	jp = &JwtProvider{}
	for _, opt := range options {
		opt(jp)
	}
	if jp.defaultTTL == 0 {
		jp.defaultTTL = time.Hour
	}
	return
}

func (jp *JwtProvider) Verify(t []byte) (verifiedToken *jwt.VerifiedToken, err error) {
	var validators []jwt.TokenValidator
	if jp.checkIssuer {
		validators = append(validators, jwt.Plain, jwt.Expected{Issuer: jp.issuer})
	}
	verifiedToken, err = jwt.Verify(jp.alg, jp.verifyKey, t, validators...)
	return
}

func (jp *JwtProvider) Gen(p JwtPayload) (t []byte, err error) {
	id := uuid.New().String()
	sub := fmt.Sprintf("%s_%s_%s", p.Dept, p.Client, p.Uid)
	standardClaims := jwt.Claims{
		NotBefore: time.Now().Unix(),
		ID:        id,
		Issuer:    jp.issuer,
		Subject:   sub,
		Audience:  []string{p.Client},
	}
	t, err = jwt.Sign(jwt.HS256, jp.signKey, p, standardClaims, jwt.MaxAge(jp.defaultTTL))
	return
}

func (jp *JwtProvider) GenWithTTL(p JwtPayload, ttl time.Duration) (t []byte, err error) {
	id := uuid.New().String()
	sub := fmt.Sprintf("%s_%s_%s", p.Dept, p.Client, p.Uid)
	standardClaims := jwt.Claims{
		NotBefore: time.Now().Unix(),
		ID:        id,
		Issuer:    jp.issuer,
		Subject:   sub,
		Audience:  []string{p.Client},
	}
	t, err = jwt.Sign(jwt.HS256, jp.signKey, p, standardClaims, jwt.MaxAge(ttl))
	return
}

func (jp *JwtProvider) Decode(t []byte) (ut *jwt.UnverifiedToken, err error) {
	ut, err = jwt.Decode(t)
	return
}
