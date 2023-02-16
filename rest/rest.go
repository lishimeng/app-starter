package rest

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type Rest interface {
	Get(uri string) (code int, body string, err error)
	GetJson(uri string, body interface{}) (code int, err error)
	Post(uri string) (code int, result string, err error)
	PostJson(uri string, req interface{}, result interface{}) (code int, err error)
}

type Handler struct {
	proxy *resty.Client
}

type ReqOption func()

func New() (r Rest) {

	h := Handler{proxy: resty.New()}
	r = &h
	return
}

func (r *Handler) Get(uri string) (code int, result string, err error) {

	resp, err := r.proxy.R().Get(uri)
	if err != nil {
		return
	}

	code = resp.StatusCode()
	bodyBs := resp.Body()
	result = string(bodyBs)
	return
}

func (r *Handler) GetJson(uri string, body interface{}) (code int, err error) {

	resp, err := r.proxy.R().Get(uri)
	if err != nil {
		return
	}

	code = resp.StatusCode()
	txt := resp.Body()
	err = json.Unmarshal(txt, body)
	return
}

func (r *Handler) Post(uri string) (code int, body string, err error) {
	resp, err := r.proxy.R().Post(uri)
	if err != nil {
		return
	}
	code = resp.StatusCode()
	body = string(resp.Body())
	return
}

func (r *Handler) PostJson(uri string, req interface{}, result interface{}) (code int, err error) {
	resp, err := r.proxy.NewRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&result).
		Post(uri)
	if err != nil {
		return
	}
	code = resp.StatusCode()
	return
}
