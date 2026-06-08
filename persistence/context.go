package persistence

type OrmContext struct {
	session Session
	// Deprecated: use facade methods (Query, Transaction). Will be removed in a future version.
	Context any
}

type TxContext struct {
	Tx Tx
	// Deprecated: use facade methods (Query, Insert, Update, Delete, Get). Will be removed in a future version.
	Context any
}

func New() *OrmContext {
	if s := resolveSession(DefaultAlias); s != nil {
		return WrapSession(s)
	}
	return &OrmContext{}
}

func NewOrm(aliasName string) *OrmContext {
	if aliasName == "" {
		aliasName = DefaultAlias
	}
	if s := resolveSession(aliasName); s != nil {
		return WrapSession(s)
	}
	return &OrmContext{}
}


func (o *OrmContext) Query(model any) Query {
	if o == nil || o.session == nil {
		return nil
	}
	return o.session.Query(model)
}

func (o *OrmContext) NewCondition() Condition {
	return NewCondition()
}

func (t *TxContext) Query(model any) Query {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Query(model)
}

func (t *TxContext) NewCondition() Condition {
	return NewCondition()
}

func (t *TxContext) Insert(model any) error {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Insert(model)
}

func (t *TxContext) Update(model any, cols ...string) error {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Update(model, cols...)
}

func (t *TxContext) Delete(model any, cols ...string) error {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Delete(model, cols...)
}

func (t *TxContext) Get(model any, cols ...string) error {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Get(model, cols...)
}

func (t *TxContext) Raw(sql string, args ...any) Query {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Raw(sql, args...)
}

func (o *OrmContext) SetLogEnable(enable bool) {
	if o == nil || o.session == nil {
		return
	}
	o.session.SetDebug(enable)
}

func (o *OrmContext) Transaction(h func(TxContext) error) (err error) {
	if o == nil || o.session == nil || h == nil {
		return nil
	}
	return o.session.Transaction(func(tx Tx) error {
		return h(WrapTx(tx))
	})
}

// WrapSession builds an OrmContext from a Session and populates legacy Context when supported.
func WrapSession(s Session) *OrmContext {
	ctx := &OrmContext{session: s}
	if exposer, ok := s.(legacyOrmExposer); ok {
		ctx.Context = exposer.LegacyOrmer()
	}
	return ctx
}

// WrapTx builds a TxContext from a Tx and populates legacy Context when supported.
func WrapTx(tx Tx) TxContext {
	ctx := TxContext{Tx: tx}
	if exposer, ok := tx.(legacyTxExposer); ok {
		ctx.Context = exposer.LegacyTxOrmer()
	}
	return ctx
}

type legacyOrmExposer interface {
	LegacyOrmer() any
}

type legacyTxExposer interface {
	LegacyTxOrmer() any
}
