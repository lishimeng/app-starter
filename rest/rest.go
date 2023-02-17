package rest

import (
	"github.com/go-resty/resty/v2"
)

type Header struct {
	Name  string
	Value string
}

type Rest interface {
	Get(uri string) (code int, body string, err error)
	GetJson(uri string, body interface{}, result interface{}, headers ...Header) (code int, err error)
	Post(uri string) (code int, result string, err error)
	PostJson(uri string, req interface{}, result interface{}, headers ...Header) (code int, err error)
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

func (r *Handler) Post(uri string) (code int, body string, err error) {
	resp, err := r.proxy.R().Post(uri)
	if err != nil {
		return
	}
	code = resp.StatusCode()
	body = string(resp.Body())
	return
}

func (r *Handler) GetJson(uri string, req interface{}, result interface{}, headers ...Header) (code int, err error) {
	client := r.proxy.NewRequest()
	for _, h := range headers {
		client = client.SetHeader(h.Name, h.Value)
	}
	resp, err := client.
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&result).
		Get(uri)
	if err != nil {
		return
	}
	code = resp.StatusCode()
	return
}

func (r *Handler) PostJson(uri string, req interface{}, result interface{}, headers ...Header) (code int, err error) {

	client := r.proxy.NewRequest()
	for _, h := range headers {
		client = client.SetHeader(h.Name, h.Value)
	}
	resp, err := client.
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
