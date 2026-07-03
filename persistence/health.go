package persistence

import (
	"context"
	"fmt"
	"time"
)

const defaultPingTimeout = 2 * time.Second

// Ping checks database connectivity for alias. No-op when the alias is not configured.
func Ping(alias string) error {
	if alias == "" {
		alias = DefaultAlias
	}
	s := resolveSession(alias)
	if s == nil {
		return nil
	}
	gs, ok := s.(*gormSession)
	if !ok || gs.db == nil {
		return nil
	}
	sqlDB, err := gs.db.DB()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()
	if err = sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("persistence: ping %q: %w", alias, err)
	}
	return nil
}
