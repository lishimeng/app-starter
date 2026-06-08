package persistence

import gormdb "gorm.io/gorm"

type gormQuery struct {
	db    *gormdb.DB
	model any
}

func wrapGormQuery(db *gormdb.DB, model any) Query {
	if db == nil {
		return nil
	}
	q := db
	if model != nil {
		q = db.Model(model)
	}
	return &gormQuery{db: q, model: model}
}

func (q *gormQuery) Filter(expr string, args ...any) Query {
	if q == nil || q.db == nil || len(args) == 0 {
		return q
	}
	if isSQLExpr(expr) {
		return &gormQuery{db: q.db.Where(expr, args...), model: q.model}
	}
	if len(args) == 1 {
		clause, clauseArgs := fieldFilter(expr, args[0])
		return &gormQuery{db: q.db.Where(clause, clauseArgs...), model: q.model}
	}
	return &gormQuery{db: q.db.Where(expr, args...), model: q.model}
}

func (q *gormQuery) FilterCond(cond Condition) Query {
	if q == nil || q.db == nil {
		return q
	}
	c, ok := cond.(*gormCondition)
	if !ok || c == nil || c.expr == "" {
		return q
	}
	return &gormQuery{db: q.db.Where(c.expr, c.args...), model: q.model}
}

func (q *gormQuery) OrderBy(expr ...string) Query {
	if q == nil || q.db == nil {
		return q
	}
	db := q.db
	for _, order := range applyOrderExprs(expr...) {
		db = db.Order(order)
	}
	return &gormQuery{db: db, model: q.model}
}

func (q *gormQuery) Offset(n int) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Offset(n), model: q.model}
}

func (q *gormQuery) Limit(n int) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Limit(n), model: q.model}
}

func (q *gormQuery) Count() (int64, error) {
	if q == nil || q.db == nil {
		return 0, nil
	}
	var count int64
	err := q.db.Count(&count).Error
	return count, err
}

func (q *gormQuery) All(dest any) (int64, error) {
	if q == nil || q.db == nil {
		return 0, nil
	}
	result := q.db.Find(dest)
	return result.RowsAffected, result.Error
}

func (q *gormQuery) One(dest any) error {
	if q == nil || q.db == nil {
		return nil
	}
	return q.db.Take(dest).Error
}

// GormDB exposes the underlying *gorm.DB for a Query when GORM-backed.
func GormDB(q Query) (*gormdb.DB, bool) {
	wrapped, ok := q.(*gormQuery)
	if !ok || wrapped == nil || wrapped.db == nil {
		return nil, false
	}
	return wrapped.db, true
}
