package digest

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestVerify(t *testing.T) {
	var plaintext = "f2383236"
	var secure = "2d2cce493b461c6d13d63f90b2696ee6d78241501f9c454d20056fd9509f120e"
	var ct, err = time.Parse("2006-01-02 15:04:05+00", "2023-02-21 12:04:25+00")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ct.Local())
	var nano = ct.UnixNano()
	ok := Verify(plaintext, secure, nano)
	if !ok {
		t.Fail()
	}
	t.Log(ok)
}

func TestVerifyAfterGenerate(t *testing.T) {
	var plaintext = "f2383236"
	var ct, err = time.Parse("2006-01-02 15:04:05+00", "2023-02-21 12:04:25+00")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ct.Local())
	var nano = ct.UnixNano()

	secure := Generate(plaintext, nano)
	ok := Verify(plaintext, secure, nano)
	if !ok {
		t.Fail()
	}
	t.Log(ok)
}

func TestVerifyAfterAlgGenerate(t *testing.T) {
	var plaintext = "f2383236"
	var ct, err = time.Parse("2006-01-02 15:04:05+00", "2023-02-21 12:04:25+00")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ct.Local())
	var nano = ct.UnixNano()

	secure := GenerateWithAlg(plaintext, nano, Digests["SM3"])
	ok := Verify(plaintext, secure, nano)
	if !ok {
		t.Fail()
	}
	t.Log(ok)
}

func TestVerifyBcryptAlgGenerate(t *testing.T) {
	var plaintext = "f2383236"
	var ct, err = time.Parse("2006-01-02 15:04:05+00", "2023-02-21 12:04:25+00")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ct.Local())
	var nano = ct.UnixNano()

	secure := GenerateWithAlg(plaintext, nano, Digests["BCRYPT"])
	ok := VerifyWithAlg(plaintext, secure, nano, Verifiers["BCRYPT"])
	if !ok {
		t.Fail()
	}
	t.Log(ok)
}

func TestBcryptAlg(t *testing.T) {
	var plaintext = "f2383236"

	var hashed, err = bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(string(hashed))
	err = bcrypt.CompareHashAndPassword(hashed, []byte(plaintext))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log("ok")
}
