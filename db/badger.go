package db

import (
	"github.com/dgraph-io/badger/v4"
)

// BadgerDB implements the DB interface using Badger
type BadgerDB struct {
	db *badger.DB
}

// NewBadgerDB creates a new BadgerDB instance
func NewBadgerDB(path string) (DB, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return &BadgerDB{db: db}, nil
}

// Get retrieves a value for the given key
func (b *BadgerDB) Get(key []byte) ([]byte, error) {
	var value []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrKeyNotFound
			}
			return err
		}

		return item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	return value, nil
}

// Set stores a key-value pair
func (b *BadgerDB) Set(key []byte, value []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// BeginTx starts a new transaction
func (b *BadgerDB) BeginTx() (Transaction, error) {
	return &BadgerTransaction{txn: b.db.NewTransaction(true)}, nil
}

// Close closes the database
func (b *BadgerDB) Close() error {
	return b.db.Close()
}

// BadgerTransaction implements the Transaction interface for Badger
type BadgerTransaction struct {
	txn *badger.Txn
}

// Set stores a key-value pair within a transaction
func (t *BadgerTransaction) Set(key []byte, value []byte) error {
	return t.txn.Set(key, value)
}

// Commit commits the transaction
func (t *BadgerTransaction) Commit() error {
	return t.txn.Commit()
}

// Rollback aborts the transaction
func (t *BadgerTransaction) Rollback() error {
	t.txn.Discard()
	return nil
}
