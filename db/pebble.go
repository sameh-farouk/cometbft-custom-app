package db

import (
	"github.com/cockroachdb/pebble"
)

// PebbleDB implements the DB interface using Pebble
type PebbleDB struct {
	db *pebble.DB
}

// NewPebbleDB creates a new PebbleDB instance
func NewPebbleDB(path string) (DB, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &PebbleDB{db: db}, nil
}

// Get retrieves a value for the given key
func (p *PebbleDB) Get(key []byte) ([]byte, error) {
	value, closer, err := p.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	defer closer.Close()

	// Copy the value since it will be invalid after closer.Close()
	result := make([]byte, len(value))
	copy(result, value)
	return result, nil
}

// Set stores a key-value pair
func (p *PebbleDB) Set(key []byte, value []byte) error {
	return p.db.Set(key, value, pebble.Sync)
}

// Delete removes a key-value pair
func (p *PebbleDB) Delete(key []byte) error {
	return p.db.Delete(key, pebble.Sync)
}

// BeginTx starts a new transaction
func (p *PebbleDB) BeginTx() (Transaction, error) {
	batch := p.db.NewBatch()
	return &PebbleTransaction{
		db:    p.db,
		batch: batch}, nil
}

// Close closes the database
func (p *PebbleDB) Close() error {
	return p.db.Close()
}

// PebbleTransaction implements the Transaction interface for Pebble
type PebbleTransaction struct {
	db    *pebble.DB
	batch *pebble.Batch
}

// Set stores a key-value pair within a transaction
func (t *PebbleTransaction) Set(key []byte, value []byte) error {
	return t.batch.Set(key, value, nil)
}

// Commit commits the transaction
func (t *PebbleTransaction) Commit() error {
	return t.batch.Commit(pebble.Sync)
}

// Rollback aborts the transaction
func (t *PebbleTransaction) Rollback() error {
	t.batch.Close()
	return nil
}
