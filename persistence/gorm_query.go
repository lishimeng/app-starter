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

func (q *gormQuery) Where(query interface{}, args ...interface{}) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Where(query, args...), model: q.model}
}

func (q *gormQuery) Or(query interface{}, args ...interface{}) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Or(query, args...), model: q.model}
}

func (q *gormQuery) Not(query interface{}, args ...interface{}) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Not(query, args...), model: q.model}
}

func (q *gormQuery) Select(query interface{}, args ...interface{}) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Select(query, args...), model: q.model}
}

func (q *gormQuery) Order(value interface{}) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Order(value), model: q.model}
}

func (q *gormQuery) Offset(offset int) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Offset(offset), model: q.model}
}

func (q *gormQuery) Limit(limit int) Query {
	if q == nil || q.db == nil {
		return q
	}
	return &gormQuery{db: q.db.Limit(limit), model: q.model}
}

func (q *gormQuery) Count() (int64, error) {
	if q == nil || q.db == nil {
		return 0, nil
	}
	var count int64
	err := q.db.Count(&count).Error
	return count, err
}

func (q *gormQuery) Find(dest interface{}, conds ...interface{}) error {
	if q == nil || q.db == nil {
		return nil
	}
	tx := q.db
	if len(conds) > 0 {
		tx = tx.Where(conds[0], conds[1:]...)
	}
	return tx.Find(dest).Error
}

func (q *gormQuery) First(dest interface{}, conds ...interface{}) error {
	if q == nil || q.db == nil {
		return nil
	}
	tx := q.db
	if len(conds) > 0 {
		tx = tx.Where(conds[0], conds[1:]...)
	}
	return tx.First(dest).Error
}

func (q *gormQuery) Take(dest interface{}, conds ...interface{}) error {
	if q == nil || q.db == nil {
		return nil
	}
	tx := q.db
	if len(conds) > 0 {
		tx = tx.Where(conds[0], conds[1:]...)
	}
	return tx.Take(dest).Error
}

func (q *gormQuery) Updates(value interface{}) error {
	if q == nil || q.db == nil {
		return nil
	}
	return q.db.Updates(value).Error
}

// GormDB exposes the underlying *gorm.DB for a Query.
func GormDB(q Query) (*gormdb.DB, bool) {
	wrapped, ok := q.(*gormQuery)
	if !ok || wrapped == nil || wrapped.db == nil {
		return nil, false
	}
	return wrapped.db, true
}
