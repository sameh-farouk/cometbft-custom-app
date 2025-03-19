package db

// DB defines the interface for database operations
type DB interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
	BeginTx() (Transaction, error)
	Close() error
}

// Transaction defines the interface for transaction operations
type Transaction interface {
	Set(key []byte, value []byte) error
	Commit() error
	Rollback() error
}
