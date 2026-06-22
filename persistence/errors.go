package persistence

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var errNotFound = errors.New("persistence: not found")

// ErrNotFound is returned by framework APIs after NormalizeErr (not gorm.ErrRecordNotFound).
var ErrNotFound = errNotFound

// NormalizeErr maps gorm driver errors to persistence-level errors.
func NormalizeErr(err error) error {
	if err == nil {
		return nil
	}
	if IsGormRecordNotFound(err) {
		return errNotFound
	}
	return err
}

// IsGormRecordNotFound reports raw GORM missing-record errors. Use for direct gorm.DB calls only.
func IsGormRecordNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}

// IsNotFound reports normalized not-found errors from Session / Query / Tx.
func IsNotFound(err error) bool {
	return err != nil && errors.Is(err, ErrNotFound)
}

// IsNotFoundAny reports not-found when err source is unknown or mixed.
func IsNotFoundAny(err error) bool {
	return IsGormRecordNotFound(err) || IsNotFound(err)
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
