package app

import "math"

type Response struct {
	Code    interface{} `json:"code,omitempty"`
	Success interface{} `json:"success,omitempty"`
	Message string      `json:"message,omitempty"`
}

type Pager struct {
	TotalPage int           `json:"totalPage"` // 总页数
	PageSize  int           `json:"pageSize"`  // 页面大小
	PageNum   int           `json:"pageNum"`   // 页号
	More      int           `json:"more"`      // 是否有下一页
	Data      []interface{} `json:"items,omitempty"`
}

func (p *Pager) Offset() int {
	return (p.PageNum - 1) * p.PageSize
}

func (p *Pager) Total(count int64) int {
	t := math.Ceil(float64(count) / float64(p.PageSize))
	totalPage := int(t)
	p.TotalPage = totalPage
	return totalPage
}

type PagerResponse struct {
	Response
	Pager
}

func (r *Response) SetCode(code int) *Response {
	r.Code = &code
	return r
}

func (r *Response) SetSuccess(success bool) *Response {
	r.Success = &success
	return r
}
