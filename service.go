package app

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/lishimeng/app-starter/persistence"
)

func Query(h func(ctx persistence.OrmContext) (err error)) (err error) {

	if h == nil {
		return
	}

	err = h(*GetOrm())
	return
}

func QueryList(req Pager,
	dataPtr interface{},
	queryHandler func(ctx persistence.OrmContext) (qs orm.QuerySeter),
	orderHandler func(persistence.OrmContext, orm.QuerySeter) (qs orm.QuerySeter),
	processHandler ...func(ctx persistence.OrmContext, dataPtr interface{}) error,

) (p Pager, err error) {

	context := GetOrm()

	p.PageNum = req.PageNum
	p.PageSize = req.PageSize

	qs := queryHandler(*context)
	total, err := qs.Count()
	if err != nil {
		return
	}

	p.TotalPage = p.Total(total)

	if orderHandler != nil {
		qs = orderHandler(*context, qs)
	}

	qs = qs.Offset(req.Offset()).Limit(req.PageSize)

	_, err = qs.All(dataPtr)
	if err != nil {
		return
	}

	for _, p := range processHandler {
		err = p(*context, dataPtr)
		if err != nil {
			break
		}
	}

	return
}

func Transaction(h func(ctx persistence.TxContext) error) (err error) {
	err = GetOrm().Transaction(h)
	return
}
