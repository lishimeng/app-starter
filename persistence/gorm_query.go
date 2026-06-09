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

// GormDB exposes the underlying *gorm.DB for a Query.
func GormDB(q Query) (*gormdb.DB, bool) {
	wrapped, ok := q.(*gormQuery)
	if !ok || wrapped == nil || wrapped.db == nil {
		return nil, false
	}
	return wrapped.db, true
}
