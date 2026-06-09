package persistence

// Install registers the GORM connector.
func Install() error {
	SetConnector(defaultGormConnector)
	SetGlobalDebugSetter(func(enable bool) {
		defaultGormConnector.setGlobalDebug(enable)
	})
	SetFallbackSessionFactory(func(alias string) Session {
		return defaultGormConnector.fallbackSession(alias)
	})
	return nil
}
