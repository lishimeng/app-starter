package persistence

import (
	"fmt"
	"sync"

	gormdb "gorm.io/gorm"
)

// DialectorOpener builds a GORM dialector from a DSN.
// Built-in drivers (postgres, mysql, sqlite) register automatically via *Config init.
type DialectorOpener func(dsn string) gormdb.Dialector

var (
	dialectorMu sync.RWMutex
	dialectors  = make(map[string]DialectorOpener)
)

// RegisterDialector registers a dialector for a driver name.
// Built-in configs register their driver automatically; use this only for custom drivers.
func RegisterDialector(driver string, opener DialectorOpener) {
	if driver == "" || opener == nil {
		return
	}
	dialectorMu.Lock()
	defer dialectorMu.Unlock()
	dialectors[driver] = opener
}

func resolveDialector(driver, dsn string) (gormdb.Dialector, error) {
	dialectorMu.RLock()
	opener, ok := dialectors[driver]
	dialectorMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("persistence: dialector for driver %q not registered", driver)
	}
	return opener(dsn), nil
}
