package beego

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/lishimeng/app-starter/persistence"
)

type session struct {
	alias string
	o     orm.Ormer
}

func newSession(alias string, o orm.Ormer) *session {
	return &session{alias: alias, o: o}
}

func (s *session) Transaction(fn func(persistence.Tx) error) error {
	if s == nil || s.o == nil || fn == nil {
		return nil
	}
	tx, err := s.o.Begin()
	if err != nil {
		return err
	}
	err = fn(&txWrapper{tx: tx})
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *session) Query(model any) persistence.Query {
	if s == nil || s.o == nil {
		return nil
	}
	return wrapQuery(s.o.QueryTable(model))
}

func (s *session) SetDebug(enable bool) {
	orm.Debug = enable
}

func (s *session) Alias() string {
	if s == nil {
		return ""
	}
	return s.alias
}

func (s *session) Ormer() orm.Ormer {
	if s == nil {
		return nil
	}
	return s.o
}

func (s *session) LegacyOrmer() any {
	return s.Ormer()
}

type txWrapper struct {
	tx orm.TxOrmer
}

func (t *txWrapper) Query(model any) persistence.Query {
	if t == nil || t.tx == nil {
		return nil
	}
	return wrapQuery(t.tx.QueryTable(model))
}

func (t *txWrapper) Insert(model any) error {
	if t == nil || t.tx == nil {
		return nil
	}
	_, err := t.tx.Insert(model)
	return persistence.CheckErr(err)
}

func (t *txWrapper) Update(model any, cols ...string) error {
	if t == nil || t.tx == nil {
		return nil
	}
	_, err := t.tx.Update(model, cols...)
	return err
}

func (t *txWrapper) Delete(model any, cols ...string) error {
	if t == nil || t.tx == nil {
		return nil
	}
	_, err := t.tx.Delete(model, cols...)
	return err
}

func (t *txWrapper) Get(model any, cols ...string) error {
	if t == nil || t.tx == nil {
		return nil
	}
	return t.tx.Read(model, cols...)
}

func (t *txWrapper) Raw(sql string, args ...any) persistence.Query {
	if t == nil || t.tx == nil {
		return nil
	}
	return wrapRawQuery(t.tx.Raw(sql, args...))
}

func (t *txWrapper) TxOrmer() orm.TxOrmer {
	if t == nil {
		return nil
	}
	return t.tx
}

func (t *txWrapper) LegacyTxOrmer() any {
	return t.TxOrmer()
}

// Ormer returns the underlying beego Ormer when ctx was created from a beego session.
func Ormer(ctx *persistence.OrmContext) (orm.Ormer, bool) {
	if ctx == nil {
		return nil, false
	}
	if ctx.Context != nil {
		o, ok := ctx.Context.(orm.Ormer)
		return o, ok
	}
	return nil, false
}

// TxOrmer returns the underlying beego TxOrmer when tx was created from a beego transaction.
func TxOrmer(ctx persistence.TxContext) (orm.TxOrmer, bool) {
	if ctx.Context != nil {
		o, ok := ctx.Context.(orm.TxOrmer)
		return o, ok
	}
	if ctx.Tx != nil {
		wrapped, ok := ctx.Tx.(*txWrapper)
		if !ok || wrapped == nil {
			return nil, false
		}
		return wrapped.TxOrmer(), true
	}
	return nil, false
}
