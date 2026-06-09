package persistence

// OpenOptions carries connection-level settings for a Connector.
type OpenOptions struct {
	Alias    string
	MaxIdle  int
	MaxOpen  int
	Debug    bool
	InitDB   bool
	Driver   string
	DSN        string
	DriverOpts any // driver-specific options set by *Config.Build()
}

// Connector opens database sessions.
type Connector interface {
	Open(opts OpenOptions) (Session, error)
	Migrate(alias string, models ...any) error
	RegisterModels(models ...any)
}

// Session is the unit of work for database access, typically one per alias.
type Session interface {
	Transaction(fn func(Tx) error) error
	Model(value interface{}) Query
	SetDebug(enable bool)
	Alias() string
}

// Tx represents a transactional database session.
type Tx interface {
	Model(value interface{}) Query
	Create(value interface{}) error
	Save(value interface{}) error
	Delete(value interface{}, conds ...interface{}) error
	First(dest interface{}, conds ...interface{}) error
	Raw(sql string, values ...interface{}) Query
}

// Query is a chainable GORM-style query builder.
type Query interface {
	Where(query interface{}, args ...interface{}) Query
	Or(query interface{}, args ...interface{}) Query
	Not(query interface{}, args ...interface{}) Query
	Select(query interface{}, args ...interface{}) Query
	Order(value interface{}) Query
	Offset(offset int) Query
	Limit(limit int) Query
	Count() (int64, error)
	Find(dest interface{}, conds ...interface{}) error
	First(dest interface{}, conds ...interface{}) error
	Take(dest interface{}, conds ...interface{}) error
	Updates(value interface{}) error
}
