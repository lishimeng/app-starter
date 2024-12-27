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

// SimpleQueryPage 单表分页查询(默认pageNo=1 pageSize=10)
func SimpleQueryPage[Model any, Dto any](pager *SimplePager[Model, Dto]) (err error) {
	if pager == nil || pager.Transform == nil || pager.QueryBuilder == nil {
		return
	}
	if pager.PageNum == 0 {
		pager.PageNum = 1
	}
	if pager.PageSize == 0 {
		pager.PageSize = 10
	}
	var count int64
	err = Transaction(func(tx persistence.TxContext) (er error) {
		q := pager.QueryBuilder(tx)
		if q == nil {
			return
		}
		query, ok := q.(orm.QuerySeter)
		if !ok {
			return
		}
		count, er = query.Count()
		if er != nil {
			return
		}
		if count == 0 {
			return
		}
		pager.TotalPage = pager.Total(count)
		if len(pager.OrderByExp) > 0 {
			query = query.OrderBy(pager.OrderByExp...)
		}
		_, er = query.Offset(pager.Offset()).Limit(pager.PageSize).All(&pager.DataSet)

		if er != nil {
			return
		}
		for _, data := range pager.DataSet {
			var dto Dto
			pager.Transform(data, &dto)
			pager.Data = append(pager.Data, dto)
		}
		return
	})
	return
}
