package persistence

import "github.com/beego/beego/v2/client/orm"

type OrmContext struct {
	Context orm.Ormer
}

func (o *OrmContext) NewCondition() *orm.Condition {
	return orm.NewCondition()
}

type TxContext struct {
	Context orm.TxOrmer
}

func (o *TxContext) NewCondition() *orm.Condition {
	return orm.NewCondition()
}

func New() *OrmContext {
	c := new(OrmContext)
	c.Context = orm.NewOrm()
	return c
}

func NewOrm(aliasName string) *OrmContext {
	c := new(OrmContext)
	c.Context = orm.NewOrmUsingDB(aliasName)
	return c
}

func (o *OrmContext) SetLogEnable(enable bool) {
	orm.Debug = enable
}

func (o *OrmContext) Transaction(h func(TxContext) error) (err error) {

	var ctx TxContext
	if h == nil {
		return
	}
	ctx.Context, err = o.Context.Begin()

	if err != nil {
		return
	}

	err = h(ctx)

	if err != nil {
		_ = ctx.Context.Rollback()
	} else {
		err = ctx.Context.Commit()
	}

	return
}
