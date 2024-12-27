package app

import (
	"github.com/lishimeng/app-starter/persistence"
	"math"
)

type Response struct {
	Code    int         `json:"code,omitempty"`
	Success string      `json:"success,omitempty"`
	Message string      `json:"message,omitempty"`
	Status  interface{} `json:"status,omitempty"`
}

type ResponseWrapper struct {
	Response
	Data any `json:"data,omitempty"`
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

type SimplePager[DbModel any, Dto any] struct {
	Pager
	DataSet      []DbModel
	Transform    func(src DbModel, dst *Dto)
	OrderByExp   []string
	QueryBuilder func(tx persistence.TxContext) any
}

type PagerResponse struct {
	Response
	Pager
}
