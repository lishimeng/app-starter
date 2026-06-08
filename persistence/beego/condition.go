package beego

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/lishimeng/app-starter/persistence"
)

type condition struct {
	c *orm.Condition
}

func newCondition() *condition {
	return &condition{c: orm.NewCondition()}
}

func (c *condition) underlying() *orm.Condition {
	if c == nil || c.c == nil {
		return orm.NewCondition()
	}
	return c.c
}

func (c *condition) And(expr string, args ...any) persistence.Condition {
	return &condition{c: c.underlying().And(expr, args...)}
}

func (c *condition) Or(expr string, args ...any) persistence.Condition {
	return &condition{c: c.underlying().Or(expr, args...)}
}

func (c *condition) AndCond(cond persistence.Condition) persistence.Condition {
	other, ok := cond.(*condition)
	if !ok || other == nil {
		return c
	}
	return &condition{c: c.underlying().AndCond(other.underlying())}
}

func (c *condition) OrCond(cond persistence.Condition) persistence.Condition {
	other, ok := cond.(*condition)
	if !ok || other == nil {
		return c
	}
	return &condition{c: c.underlying().OrCond(other.underlying())}
}
