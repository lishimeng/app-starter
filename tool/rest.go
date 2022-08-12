package tool

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type Rest interface {
}

type RestHandler struct {
	proxy *resty.Client
}

func New() (r Rest) {

	h := RestHandler{proxy: resty.New()}
	r = &h
	return
}

func (r *RestHandler) Get(uri string) (code int, body string, err error) {

	resp, err := r.proxy.R().Get(uri)
	if err != nil {
		return
	}

	code = resp.StatusCode()
	bodyBs := resp.Body()
	body = string(bodyBs)
	return
}

func (r *RestHandler) GetJson(uri string, body interface{}) (code int, err error) {

	resp, err := r.proxy.R().Get(uri)
	if err != nil {
		return
	}

	code = resp.StatusCode()
	txt := resp.Body()
	err = json.Unmarshal(txt, body)
	return
}

func (r *RestHandler) Post(uri string) (code int, body string, err error) {
	resp, err := r.proxy.R().Post(uri)
	if err != nil {
		return
	}
	code = resp.StatusCode()
	body = string(resp.Body())
	return
}

func (r *RestHandler) PostJson(uri string, body interface{}) (code int, err error) {
	resp, err := r.proxy.R().Post(uri)
	if err != nil {
		return
	}
	code = resp.StatusCode()
	txt := resp.Body()
	err = json.Unmarshal(txt, body)
	return
}
