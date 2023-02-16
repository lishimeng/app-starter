package token

import (
	"testing"
	"time"
)

func TestJwt001(t *testing.T) {

	var payload = JwtPayload{
		Org:    "domain.com",
		Uid:    "wx_23432fve4",
		Client: "wechat",
	}
	var provider = NewJwtProvider("www.demo.com",
		WithAlg("HS256"),
		WithKey(sharedKey, sharedKey),
		WithDefaultTTL(time.Hour*2))
	var bs, err = provider.Gen(payload)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf(string(bs))

	vt, err := provider.Verify(bs)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%+v\n", vt.StandardClaims)
	t.Log(string(vt.Payload))

	var p JwtPayload
	err = vt.Claims(&p)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%+v\n", p)
}

func TestJwtDecode(t *testing.T) {

	var payload = JwtPayload{
		Org:    "domain.com",
		Uid:    "wx_23432fve4",
		Client: "wechat",
	}
	var provider = NewJwtProvider("www.demo.com",
		WithAlg("HS256"),
		WithKey(sharedKey, nil),
		WithDefaultTTL(time.Hour*2))
	var bs, err = provider.Gen(payload)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf(string(bs))

	vt, err := provider.Decode(bs)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(vt.Payload))

	var p JwtPayload
	err = vt.Claims(&p)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%+v\n", p)
}
