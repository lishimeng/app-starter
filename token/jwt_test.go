package token

import (
	"testing"
	"time"
)

func TestJwtProviderSignVerify(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	prov := NewJwtProvider(
		WithAlg("HS256"),
		WithKey(key, key),
		WithIssuer("test-issuer"),
		WithDefaultTTL(time.Hour),
	)

	payload := JwtPayload{
		Uid:    "u1",
		Client: "pc",
		Dept:   "d1",
		Scope:  "read",
	}
	raw, err := prov.Gen(payload)
	if err != nil {
		t.Fatal(err)
	}

	verified, err := prov.Verify(raw)
	if err != nil {
		t.Fatal(err)
	}
	var got JwtPayload
	if err = verified.Claims(&got); err != nil {
		t.Fatal(err)
	}
	if got.Uid != payload.Uid || got.Client != payload.Client {
		t.Fatalf("payload mismatch: %+v", got)
	}
}

func TestDecodeUnverified(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	prov := NewJwtProvider(WithKey(key, key))
	raw, err := prov.Gen(JwtPayload{Uid: "u2", Client: "app"})
	if err != nil {
		t.Fatal(err)
	}

	ut, err := Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	var got JwtPayload
	if err = ut.Claims(&got); err != nil {
		t.Fatal(err)
	}
	if got.Uid != "u2" {
		t.Fatalf("uid = %q, want u2", got.Uid)
	}
}

func TestVerifyWrongIssuer(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	prov := NewJwtProvider(
		WithKey(key, key),
		WithIssuer("expected"),
	)
	raw, err := prov.Gen(JwtPayload{Uid: "u3"})
	if err != nil {
		t.Fatal(err)
	}

	checker := NewJwtProvider(
		WithKey(key, key),
		WithIssuer("other"),
	)
	if _, err = checker.Verify(raw); err == nil {
		t.Fatal("expected issuer mismatch error")
	}
}
