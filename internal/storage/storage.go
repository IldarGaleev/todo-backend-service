package storage

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAccessDenied  = errors.New("access denied")
	ErrDatabaseError = errors.New("database error")
)
