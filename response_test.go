package app

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestResponse001(t *testing.T) {
	r := PagerResponse{
		Response{},
		Pager{

		},
	}
	r.Response.SetCode(1)
	r.Response.Message="success"
	bs, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Printf(string(bs))
}

func TestResponse002(t *testing.T) {
	r := PagerResponse{
		Response: Response{},
	}
	r.Response.SetCode(1)
	r.Response.Message="success"
	bs, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Printf(string(bs))
}

func TestResponse003(t *testing.T) {
	r := Response{
	}
	r.SetCode(1)
	r.Message="success"
	bs, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Printf(string(bs))
}