package persistence

import "sync"

const DefaultAlias = "default"

var (
	connectorMu sync.RWMutex
	connector   Connector

	sessionMu sync.RWMutex
	sessions  = make(map[string]Session)

	conditionMu      sync.RWMutex
	conditionFactory func() Condition

	fallbackMu            sync.RWMutex
	fallbackSessionFactory func(alias string) Session
)

// SetConnector installs the active database connector. Called from Install().
func SetConnector(c Connector) {
	connectorMu.Lock()
	defer connectorMu.Unlock()
	connector = c
}

func getConnector() Connector {
	connectorMu.RLock()
	defer connectorMu.RUnlock()
	return connector
}

// SetConditionFactory registers how to create query conditions.
func SetConditionFactory(factory func() Condition) {
	conditionMu.Lock()
	defer conditionMu.Unlock()
	conditionFactory = factory
}

// NewCondition creates a condition instance.
func NewCondition() Condition {
	conditionMu.RLock()
	factory := conditionFactory
	conditionMu.RUnlock()
	if factory == nil {
		return nil
	}
	return factory()
}

// RegisterSession stores a session for the given alias.
func RegisterSession(alias string, session Session) {
	if alias == "" {
		alias = DefaultAlias
	}
	sessionMu.Lock()
	defer sessionMu.Unlock()
	sessions[alias] = session
}

// GetSession returns the session registered for alias, or nil.
func GetSession(alias string) Session {
	if alias == "" {
		alias = DefaultAlias
	}
	sessionMu.RLock()
	defer sessionMu.RUnlock()
	return sessions[alias]
}

// SetFallbackSessionFactory provides sessions when no alias has been registered yet.
func SetFallbackSessionFactory(factory func(alias string) Session) {
	fallbackMu.Lock()
	defer fallbackMu.Unlock()
	fallbackSessionFactory = factory
}

func resolveSession(alias string) Session {
	if s := GetSession(alias); s != nil {
		return s
	}
	fallbackMu.RLock()
	factory := fallbackSessionFactory
	fallbackMu.RUnlock()
	if factory == nil {
		return nil
	}
	return factory(alias)
}
