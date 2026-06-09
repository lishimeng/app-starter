package persistence

import (
	"fmt"
	"strings"
)

func validatePostgresDSN(dsn string) error {
	required := map[string]bool{"user": false, "dbname": false, "host": false}
	for _, part := range strings.Fields(dsn) {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		if _, exists := required[key]; exists && value != "" {
			required[key] = true
		}
	}

	var missing []string
	for key, ok := range required {
		if !ok {
			missing = append(missing, key)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf(
		"postgres config incomplete (missing %s): set DB_USER, DB_PASSWORD, DB_HOST, DB_DATABASE environment variables",
		strings.Join(missing, ", "),
	)
}
