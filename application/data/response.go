package web

type Response struct {
	Code *int `json:"code,omitempty"`
	Success *bool `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

type Pager struct {
	TotalPage int `json:"totalPage"`// 总页数
	PageSize int `json:"pageSize"`// 页面大小
	PageNum int `json:"pageNum"` // 页号
	More int `json:"more"` // 是否有下一页
	Data []interface{} `json:"items,omitempty"`
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
