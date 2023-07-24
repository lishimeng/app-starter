package miniwx

import "testing"

func TestAuthorizeCode(t *testing.T) {

	handler := New(WithAuth("", ""))
	resp, err := handler.AuthorizeCode("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("login success!")
	t.Log(resp.OpenId)
	t.Log(resp.UnionId)
	t.Log(resp.SessionKey)
}

func TestCredentials(t *testing.T) {

	handler := New(WithAuth("", "")) // no cache
	t.Log("get new credential:")
	resp, err := handler.Credentials()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("get credential success!")
	t.Log(resp)

	t.Log("get credential again:")
	resp, err = handler.Credentials()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("get the next credential success!")
	t.Log(resp)
}
