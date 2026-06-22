package persistence

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var errNotFound = errors.New("persistence: not found")

// NormalizeErr maps gorm driver errors to persistence-level errors.
func NormalizeErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errNotFound
	}
	return err
}

// IsNotFound reports whether err is a missing-record error from persistence.
func IsNotFound(err error) bool {
	return errors.Is(err, errNotFound)
}

// IsDuplicate reports common unique-constraint violations.
func IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "duplicate key")
}
