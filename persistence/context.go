package persistence

type OrmContext struct {
	session Session
	// Deprecated: use facade methods. Exposes *gorm.DB when available.
	Context any
}

type TxContext struct {
	Tx Tx
	// Deprecated: use facade methods. Exposes *gorm.DB when available.
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

func (o *OrmContext) Model(value interface{}) Query {
	if o == nil || o.session == nil {
		return nil
	}
	return o.session.Model(value)
}

func (t *TxContext) Model(value interface{}) Query {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Model(value)
}

func (t *TxContext) Create(value interface{}) error {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Create(value)
}

func (t *TxContext) Save(value interface{}) error {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Save(value)
}

func (t *TxContext) Delete(value interface{}, conds ...interface{}) error {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Delete(value, conds...)
}

func (t *TxContext) First(dest interface{}, conds ...interface{}) error {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.First(dest, conds...)
}

func (t *TxContext) Raw(sql string, values ...interface{}) Query {
	if t == nil || t.Tx == nil {
		return nil
	}
	return t.Tx.Raw(sql, values...)
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

func WrapSession(s Session) *OrmContext {
	ctx := &OrmContext{session: s}
	if exposer, ok := s.(legacyOrmExposer); ok {
		ctx.Context = exposer.LegacyOrmer()
	}
	return ctx
}

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
