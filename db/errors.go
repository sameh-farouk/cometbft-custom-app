package db

import (
	"errors"
)

// Standard errors
var (
	ErrKeyNotFound = errors.New("key not found")
	ErrTxnConflict = errors.New("transaction conflict")
	ErrDBClosed    = errors.New("database is closed")
)
