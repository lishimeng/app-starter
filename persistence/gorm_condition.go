package persistence

type gormCondition struct {
	expr string
	args []any
}

func newGormCondition() *gormCondition {
	return &gormCondition{}
}

func (c *gormCondition) merge(op, expr string, args ...any) *gormCondition {
	if c == nil {
		return &gormCondition{expr: expr, args: append([]any(nil), args...)}
	}
	if c.expr == "" {
		c.expr = expr
		c.args = append(c.args, args...)
		return c
	}
	c.expr = "(" + c.expr + ") " + op + " (" + expr + ")"
	c.args = append(c.args, args...)
	return c
}

func (c *gormCondition) And(expr string, args ...any) Condition {
	return c.merge("AND", expr, args...)
}

func (c *gormCondition) Or(expr string, args ...any) Condition {
	return c.merge("OR", expr, args...)
}

func (c *gormCondition) AndCond(cond Condition) Condition {
	other, ok := cond.(*gormCondition)
	if !ok || other == nil || other.expr == "" {
		return c
	}
	return c.merge("AND", other.expr, other.args...)
}

func (c *gormCondition) OrCond(cond Condition) Condition {
	other, ok := cond.(*gormCondition)
	if !ok || other == nil || other.expr == "" {
		return c
	}
	return c.merge("OR", other.expr, other.args...)
}
