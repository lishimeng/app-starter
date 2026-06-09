package persistence

import (
	"fmt"
	"sync"

	gormdb "gorm.io/gorm"
)

// DialectorOpener builds a GORM dialector from OpenOptions.
// Built-in drivers (postgres, mysql, sqlite) register automatically via *Config init.
type DialectorOpener func(opts OpenOptions) gormdb.Dialector

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

func resolveDialector(opts OpenOptions) (gormdb.Dialector, error) {
	dialectorMu.RLock()
	opener, ok := dialectors[opts.Driver]
	dialectorMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("persistence: dialector for driver %q not registered", opts.Driver)
	}
	return opener(opts), nil
}
