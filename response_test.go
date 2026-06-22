package app

import (
	"encoding/json"
	"testing"
)

func TestResponse001(t *testing.T) {
	r := PagerResponse{
		Response{},
	}
	r.Response.Code = 1
	r.Response.Message = "success"
	bs, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(bs))
}

func TestResponse002(t *testing.T) {
	r := PagerResponse{
		Response: Response{},
	}
	r.Response.Code = 1
	r.Response.Message = "success"
	bs, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(bs))
}

func TestResponse003(t *testing.T) {
	r := Response{}
	r.Code = 1
	r.Message = "success"
	bs, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(bs))
}
