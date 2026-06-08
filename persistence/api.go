package persistence

// OpenOptions carries connection-level settings for a Connector.
type OpenOptions struct {
	Alias    string
	MaxIdle  int
	MaxOpen  int
	Debug    bool
	InitDB   bool
	Driver   string
	DSN      string
	DBParams []any // implementation-specific pool/connection options
}

// Connector opens database sessions. Implementations live outside the public SDK surface
// (e.g. persistence/beego, persistence/gorm).
type Connector interface {
	Open(opts OpenOptions) (Session, error)
	Migrate(alias string, models ...any) error
	RegisterModels(models ...any)
}

// Session is the unit of work for database access, typically one per alias.
type Session interface {
	Transaction(fn func(Tx) error) error
	Query(model any) Query
	SetDebug(enable bool)
	Alias() string
}

// Tx represents a transactional database session.
type Tx interface {
	Query(model any) Query
	Insert(model any) error
	Update(model any, cols ...string) error
	Delete(model any, cols ...string) error
	Get(model any, cols ...string) error
	Raw(sql string, args ...any) Query
}

// Query is a chainable query builder used by QueryPage and business code.
type Query interface {
	Filter(expr string, args ...any) Query
	FilterCond(cond Condition) Query
	OrderBy(expr ...string) Query
	Offset(n int) Query
	Limit(n int) Query
	Count() (int64, error)
	All(dest any) (int64, error)
	One(dest any) error
}

// Condition composes query predicates.
type Condition interface {
	And(expr string, args ...any) Condition
	Or(expr string, args ...any) Condition
	AndCond(cond Condition) Condition
	OrCond(cond Condition) Condition
}
