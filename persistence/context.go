package persistence

type OrmContext struct {
	session Session
}

type TxContext struct {
	Tx Tx
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
	return &OrmContext{session: s}
}

func WrapTx(tx Tx) TxContext {
	return TxContext{Tx: tx}
}
