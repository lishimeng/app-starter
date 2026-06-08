package persistence

import "sync"

var (
	debugMu             sync.RWMutex
	globalDebug         bool
	globalDebugSetter   func(bool)
)

// SetGlobalDebugSetter registers a global debug toggle.
func SetGlobalDebugSetter(fn func(bool)) {
	debugMu.Lock()
	defer debugMu.Unlock()
	globalDebugSetter = fn
}

// SetDebug enables or disables SQL debug logging for all active sessions
// and the global debug switch when supported.
func SetDebug(enable bool) {
	debugMu.Lock()
	globalDebug = enable
	setter := globalDebugSetter
	debugMu.Unlock()

	if setter != nil {
		setter(enable)
	}

	sessionMu.RLock()
	defer sessionMu.RUnlock()
	for _, s := range sessions {
		s.SetDebug(enable)
	}
}

func isDebugEnabled() bool {
	debugMu.RLock()
	defer debugMu.RUnlock()
	return globalDebug
}
