package persistence

import (
	"fmt"
	"sync"

	gormdb "gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type gormConnector struct {
	mu     sync.RWMutex
	dbs    map[string]*gormdb.DB
	models []any
	debug  bool
}

func newGormConnector() *gormConnector {
	return &gormConnector{
		dbs: make(map[string]*gormdb.DB),
	}
}

var defaultGormConnector = newGormConnector()

func (c *gormConnector) Open(opts OpenOptions) (Session, error) {
	if c == nil {
		return nil, fmt.Errorf("persistence: connector nil")
	}
	alias := opts.Alias
	if alias == "" {
		alias = DefaultAlias
	}

	dialector, err := resolveDialector(opts)
	if err != nil {
		return nil, err
	}

	logLevel := logger.Silent
	if opts.Debug || c.isDebug() {
		logLevel = logger.Info
	}

	db, err := gormdb.Open(dialector, &gormdb.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if opts.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(opts.MaxIdle)
	}
	if opts.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(opts.MaxOpen)
	}

	c.mu.Lock()
	c.dbs[alias] = db
	c.mu.Unlock()

	return newGormSession(alias, db), nil
}

func (c *gormConnector) Migrate(alias string, models ...any) error {
	if c == nil {
		return fmt.Errorf("persistence: connector nil")
	}
	if alias == "" {
		alias = DefaultAlias
	}
	db, err := c.db(alias)
	if err != nil {
		return err
	}
	targets := models
	if len(targets) == 0 {
		targets = c.registeredModels()
	}
	if len(targets) == 0 {
		return nil
	}
	return db.AutoMigrate(targets...)
}

func (c *gormConnector) RegisterModels(models ...any) {
	if c == nil || len(models) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.models = append(c.models, models...)
}

func (c *gormConnector) registeredModels() []any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]any, len(c.models))
	copy(out, c.models)
	return out
}

func (c *gormConnector) db(alias string) (*gormdb.DB, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	db, ok := c.dbs[alias]
	if !ok || db == nil {
		return nil, fmt.Errorf("persistence: database alias %q not opened", alias)
	}
	return db, nil
}

func (c *gormConnector) setGlobalDebug(enable bool) {
	c.mu.Lock()
	c.debug = enable
	dbs := make([]*gormdb.DB, 0, len(c.dbs))
	for _, db := range c.dbs {
		dbs = append(dbs, db)
	}
	c.mu.Unlock()

	lvl := logger.Silent
	if enable {
		lvl = logger.Info
	}
	for _, db := range dbs {
		db.Logger = db.Logger.LogMode(lvl)
	}
}

func (c *gormConnector) isDebug() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.debug
}

func (c *gormConnector) fallbackSession(alias string) Session {
	if alias == "" {
		alias = DefaultAlias
	}
	db, err := c.db(alias)
	if err != nil {
		return nil
	}
	return newGormSession(alias, db)
}
