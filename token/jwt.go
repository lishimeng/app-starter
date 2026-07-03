package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	case App, Pad, PC, WeChat:
		jp.Client = string(clientType)
	}
	return jp
}

type jwtClaims struct {
	JwtPayload
	jwt.RegisteredClaims
}

// VerifiedToken holds cryptographically verified JWT claims.
type VerifiedToken struct {
	claims jwtClaims
}

func (vt *VerifiedToken) Claims(v interface{}) error {
	p, ok := v.(*JwtPayload)
	if !ok {
		return fmt.Errorf("token: expected *JwtPayload")
	}
	*p = vt.claims.JwtPayload
	return nil
}

// UnverifiedToken holds parsed JWT claims without signature verification.
type UnverifiedToken struct {
	claims jwtClaims
}

func (ut *UnverifiedToken) Claims(v interface{}) error {
	p, ok := v.(*JwtPayload)
	if !ok {
		return fmt.Errorf("token: expected *JwtPayload")
	}
	*p = ut.claims.JwtPayload
	return nil
}

// Decode parses a JWT without verifying its signature.
func Decode(t []byte) (*UnverifiedToken, error) {
	var c jwtClaims
	_, _, err := jwt.NewParser(jwt.WithoutClaimsValidation()).ParseUnverified(string(t), &c)
	if err != nil {
		return nil, err
	}
	return &UnverifiedToken{claims: c}, nil
}

type JwtProvider struct {
	signingMethod jwt.SigningMethod
	signKey       interface{}
	verifyKey     interface{}
	issuer        string
	defaultTTL    time.Duration
	checkIssuer   bool
}

type JwtBuildOption func(provider *JwtProvider)

var jwtAllAlg = map[string]jwt.SigningMethod{
	"HS256": jwt.SigningMethodHS256,
	"HS384": jwt.SigningMethodHS384,
	"HS512": jwt.SigningMethodHS512,
	"RS256": jwt.SigningMethodRS256,
	"RS384": jwt.SigningMethodRS384,
	"RS512": jwt.SigningMethodRS512,
	"PS256": jwt.SigningMethodPS256,
	"PS384": jwt.SigningMethodPS384,
	"PS512": jwt.SigningMethodPS512,
	"ES256": jwt.SigningMethodES256,
	"ES384": jwt.SigningMethodES384,
	"ES512": jwt.SigningMethodES512,
}

var WithAlg = func(alg string) JwtBuildOption {
	return func(provider *JwtProvider) {
		if m, ok := jwtAllAlg[alg]; ok {
			provider.signingMethod = m
		}
	}
}

// WithIssuer sets issuer; verification checks issuer when set.
var WithIssuer = func(issuer string) JwtBuildOption {
	return func(provider *JwtProvider) {
		provider.issuer = issuer
		provider.checkIssuer = true
	}
}

// WithKey sets signing and verification keys.
var WithKey = func(signKey interface{}, verifyKey interface{}) JwtBuildOption {
	return func(provider *JwtProvider) {
		provider.signKey = signKey
		provider.verifyKey = verifyKey
	}
}

// WithDefaultTTL sets the default token TTL.
var WithDefaultTTL = func(ttl time.Duration) JwtBuildOption {
	return func(provider *JwtProvider) {
		provider.defaultTTL = ttl
	}
}

func NewJwtProvider(options ...JwtBuildOption) (jp *JwtProvider) {
	jp = &JwtProvider{
		signingMethod: jwt.SigningMethodHS256,
		defaultTTL:    time.Hour,
	}
	for _, opt := range options {
		opt(jp)
	}
	return
}

func (jp *JwtProvider) Verify(t []byte) (verifiedToken *VerifiedToken, err error) {
	if jp.verifyKey == nil || jp.signingMethod == nil {
		return nil, errors.New("token: verify not configured")
	}
	var c jwtClaims
	opts := []jwt.ParserOption{}
	if jp.checkIssuer {
		opts = append(opts, jwt.WithIssuer(jp.issuer))
	}
	_, err = jwt.ParseWithClaims(string(t), &c, func(token *jwt.Token) (interface{}, error) {
		if token.Method == nil || token.Method.Alg() != jp.signingMethod.Alg() {
			return nil, fmt.Errorf("token: unexpected signing method")
		}
		return jp.verifyKey, nil
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &VerifiedToken{claims: c}, nil
}

func (jp *JwtProvider) Gen(p JwtPayload) (t []byte, err error) {
	return jp.GenWithTTL(p, jp.defaultTTL)
}

func (jp *JwtProvider) GenWithTTL(p JwtPayload, ttl time.Duration) (t []byte, err error) {
	if jp.signKey == nil || jp.signingMethod == nil {
		return nil, errors.New("token: signing not configured")
	}
	now := time.Now()
	c := jwtClaims{
		JwtPayload: p,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    jp.issuer,
			Subject:   fmt.Sprintf("%s_%s_%s", p.Dept, p.Client, p.Uid),
			Audience:  jwt.ClaimStrings{p.Client},
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	s, err := jwt.NewWithClaims(jp.signingMethod, c).SignedString(jp.signKey)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func (jp *JwtProvider) Decode(t []byte) (ut *UnverifiedToken, err error) {
	return Decode(t)
}
