package persistence

import (
	gormdb "gorm.io/gorm"
)

type gormSession struct {
	alias string
	db    *gormdb.DB
}

func newGormSession(alias string, db *gormdb.DB) *gormSession {
	return &gormSession{alias: alias, db: db}
}

func (s *gormSession) Transaction(fn func(Tx) error) error {
	if s == nil || s.db == nil || fn == nil {
		return nil
	}
	return s.db.Transaction(func(tx *gormdb.DB) error {
		return fn(&gormTx{db: tx})
	})
}

func (s *gormSession) Model(value interface{}) Query {
	if s == nil || s.db == nil {
		return nil
	}
	return wrapGormQuery(s.db, value)
}

func (s *gormSession) SetDebug(enable bool) {
	if s == nil || s.db == nil {
		return
	}
	s.db.Logger = newSlogGormLogger(gormLogLevel(enable))
}

func (s *gormSession) Alias() string {
	if s == nil {
		return ""
	}
	return s.alias
}

type gormTx struct {
	db *gormdb.DB
}

func (t *gormTx) Model(value interface{}) Query {
	if t == nil || t.db == nil {
		return nil
	}
	return wrapGormQuery(t.db, value)
}

func (t *gormTx) Create(value interface{}) error {
	if t == nil || t.db == nil {
		return nil
	}
	return NormalizeErr(CheckErr(t.db.Create(value).Error))
}

func (t *gormTx) Save(value interface{}) error {
	if t == nil || t.db == nil {
		return nil
	}
	return NormalizeErr(t.db.Save(value).Error)
}

func (t *gormTx) Delete(value interface{}, conds ...interface{}) error {
	if t == nil || t.db == nil {
		return nil
	}
	tx := t.db
	if len(conds) > 0 {
		tx = tx.Where(conds[0], conds[1:]...)
	}
	return NormalizeErr(tx.Delete(value).Error)
}

func (t *gormTx) First(dest interface{}, conds ...interface{}) error {
	if t == nil || t.db == nil {
		return nil
	}
	tx := t.db
	if len(conds) > 0 {
		tx = tx.Where(conds[0], conds[1:]...)
	}
	return NormalizeErr(tx.First(dest).Error)
}

func (t *gormTx) Raw(sql string, values ...interface{}) Query {
	if t == nil || t.db == nil {
		return nil
	}
	return wrapGormQuery(t.db.Raw(sql, values...), nil)
}

// SessionDB returns the underlying *gorm.DB for an OrmContext.
func SessionDB(ctx *OrmContext) (*gormdb.DB, bool) {
	if ctx == nil || ctx.session == nil {
		return nil, false
	}
	s, ok := ctx.session.(*gormSession)
	if !ok || s == nil || s.db == nil {
		return nil, false
	}
	return s.db, true
}

// TxDB returns the underlying transactional *gorm.DB.
func TxDB(ctx TxContext) (*gormdb.DB, bool) {
	if ctx.Tx == nil {
		return nil, false
	}
	t, ok := ctx.Tx.(*gormTx)
	if !ok || t == nil || t.db == nil {
		return nil, false
	}
	return t.db, true
}
