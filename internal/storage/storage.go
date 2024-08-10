// Package storage contains storage layer types
package storage

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("storage: not found")
	ErrAccessDenied  = errors.New("storage: access denied")
	ErrDatabaseError = errors.New("storage: database error")
)
