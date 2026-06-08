package persistence

import (
	gormdb "gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func (s *gormSession) Query(model any) Query {
	if s == nil || s.db == nil {
		return nil
	}
	return wrapGormQuery(s.db, model)
}

func (s *gormSession) SetDebug(enable bool) {
	if s == nil || s.db == nil {
		return
	}
	lvl := logger.Silent
	if enable {
		lvl = logger.Info
	}
	s.db.Logger = s.db.Logger.LogMode(lvl)
}

func (s *gormSession) Alias() string {
	if s == nil {
		return ""
	}
	return s.alias
}

func (s *gormSession) LegacyOrmer() any {
	if s == nil {
		return nil
	}
	return s.db
}

type gormTx struct {
	db *gormdb.DB
}

func (t *gormTx) Query(model any) Query {
	if t == nil || t.db == nil {
		return nil
	}
	return wrapGormQuery(t.db, model)
}

func (t *gormTx) Insert(model any) error {
	if t == nil || t.db == nil {
		return nil
	}
	return CheckErr(t.db.Create(model).Error)
}

func (t *gormTx) Update(model any, cols ...string) error {
	if t == nil || t.db == nil {
		return nil
	}
	tx := t.db
	if len(cols) > 0 {
		tx = tx.Select(cols)
	}
	return tx.Updates(model).Error
}

func (t *gormTx) Delete(model any, cols ...string) error {
	if t == nil || t.db == nil {
		return nil
	}
	tx := t.db
	if len(cols) > 0 {
		tx = tx.Select(cols)
	}
	return tx.Delete(model).Error
}

func (t *gormTx) Get(model any, cols ...string) error {
	if t == nil || t.db == nil {
		return nil
	}
	tx := t.db
	if len(cols) > 0 {
		tx = tx.Select(cols)
	}
	return tx.First(model).Error
}

func (t *gormTx) Raw(sql string, args ...any) Query {
	if t == nil || t.db == nil {
		return nil
	}
	return wrapGormQuery(t.db.Raw(sql, args...), nil)
}

func (t *gormTx) LegacyTxOrmer() any {
	if t == nil {
		return nil
	}
	return t.db
}

// SessionDB returns the underlying *gorm.DB for an OrmContext backed by GORM.
func SessionDB(ctx *OrmContext) (*gormdb.DB, bool) {
	if ctx == nil {
		return nil, false
	}
	if ctx.Context != nil {
		db, ok := ctx.Context.(*gormdb.DB)
		return db, ok
	}
	return nil, false
}

// TxDB returns the underlying transactional *gorm.DB when tx is GORM-backed.
func TxDB(ctx TxContext) (*gormdb.DB, bool) {
	if ctx.Context != nil {
		db, ok := ctx.Context.(*gormdb.DB)
		return db, ok
	}
	if ctx.Tx != nil {
		wrapped, ok := ctx.Tx.(*gormTx)
		if !ok || wrapped == nil {
			return nil, false
		}
		return wrapped.db, true
	}
	return nil, false
}
