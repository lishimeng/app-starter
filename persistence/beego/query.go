package beego

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/lishimeng/app-starter/persistence"
)

type query struct {
	q orm.QuerySeter
}

type rawQuery struct {
	q orm.RawSeter
}

func wrapQuery(q orm.QuerySeter) persistence.Query {
	if q == nil {
		return nil
	}
	return &query{q: q}
}

func wrapRawQuery(q orm.RawSeter) persistence.Query {
	if q == nil {
		return nil
	}
	return &rawQuery{q: q}
}

func (q *query) Filter(expr string, args ...any) persistence.Query {
	if q == nil || q.q == nil {
		return q
	}
	return &query{q: q.q.Filter(expr, args...)}
}

func (q *query) FilterCond(cond persistence.Condition) persistence.Query {
	if q == nil || q.q == nil {
		return q
	}
	c, ok := cond.(*condition)
	if !ok || c == nil {
		return q
	}
	return &query{q: q.q.SetCond(c.underlying())}
}

func (q *query) OrderBy(expr ...string) persistence.Query {
	if q == nil || q.q == nil {
		return q
	}
	return &query{q: q.q.OrderBy(expr...)}
}

func (q *query) Offset(n int) persistence.Query {
	if q == nil || q.q == nil {
		return q
	}
	return &query{q: q.q.Offset(n)}
}

func (q *query) Limit(n int) persistence.Query {
	if q == nil || q.q == nil {
		return q
	}
	return &query{q: q.q.Limit(n)}
}

func (q *query) Count() (int64, error) {
	if q == nil || q.q == nil {
		return 0, nil
	}
	return q.q.Count()
}

func (q *query) All(dest any) (int64, error) {
	if q == nil || q.q == nil {
		return 0, nil
	}
	return q.q.All(dest)
}

func (q *query) One(dest any) error {
	if q == nil || q.q == nil {
		return nil
	}
	return q.q.One(dest)
}

func (q *rawQuery) Filter(expr string, args ...any) persistence.Query { return q }
func (q *rawQuery) FilterCond(cond persistence.Condition) persistence.Query {
	return q
}
func (q *rawQuery) OrderBy(expr ...string) persistence.Query { return q }
func (q *rawQuery) Offset(n int) persistence.Query             { return q }
func (q *rawQuery) Limit(n int) persistence.Query              { return q }

func (q *rawQuery) Count() (int64, error) {
	return 0, fmt.Errorf("persistence/beego: Count is not supported on raw queries; use legacy TxOrmer")
}

func (q *rawQuery) All(dest any) (int64, error) {
	return 0, fmt.Errorf("persistence/beego: All is not supported on raw queries; use legacy TxOrmer")
}

func (q *rawQuery) One(dest any) error {
	return fmt.Errorf("persistence/beego: One is not supported on raw queries; use legacy TxOrmer")
}

// QuerySeter exposes the underlying beego QuerySeter for migration compatibility.
func QuerySeter(q persistence.Query) (orm.QuerySeter, bool) {
	wrapped, ok := q.(*query)
	if !ok || wrapped == nil || wrapped.q == nil {
		return nil, false
	}
	return wrapped.q, true
}
