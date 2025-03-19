//go:build tigerbeetle

package db

import (
	"fmt"
	"strconv"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// TigerBeetleDB implements the DB interface using TigerBeetle
type TigerBeetleDB struct {
	client tb.Client
}

// NewTigerBeetleDB creates a new TigerBeetleDB instance
func NewTigerBeetleDB(addresses []string) (DB, error) {
	// Convert clusterID to Uint128
	clusterIDUint128 := types.ToUint128(0)

	client, err := tb.NewClient(clusterIDUint128, addresses)
	if err != nil {
		return nil, err
	}

	return &TigerBeetleDB{
		client: client,
	}, nil
}

// isAccountKey determines if a key represents an account
func isAccountKey(key []byte) bool {
	// In our application, account keys are simple strings like "1", "2", etc.
	_, err := strconv.ParseUint(string(key), 10, 64)
	return err == nil
}

// parseAccountID converts a key to a TigerBeetle account ID
func parseAccountID(key []byte) types.Uint128 {
	id, _ := strconv.ParseUint(string(key), 10, 64)
	return types.ToUint128(id)
}

// parseBalance converts a byte slice to a balance
func parseBalance(value []byte) uint64 {
	if value == nil {
		return 0
	}
	balance, _ := strconv.ParseUint(string(value), 10, 64)
	return balance
}

// encodeBalance converts a balance to a byte slice
func encodeBalance(balance uint64) []byte {
	return []byte(strconv.FormatUint(balance, 10))
}

// getAccountBalance calculates the balance from an account
func getAccountBalance(account types.Account) uint64 {
	// In TigerBeetle, we'll use UserData64 to store the balance
	return account.UserData64
}

// Get retrieves a value for the given key
func (t *TigerBeetleDB) Get(key []byte) ([]byte, error) {
	if isAccountKey(key) {
		accountID := parseAccountID(key)
		accounts, err := t.client.LookupAccounts([]types.Uint128{accountID})
		if err != nil {
			return nil, err
		}

		if len(accounts) == 0 {
			return nil, ErrKeyNotFound
		}

		balance := getAccountBalance(accounts[0])
		return encodeBalance(balance), nil
	}

	return nil, ErrKeyNotFound
}

// Set stores a key-value pair
func (t *TigerBeetleDB) Set(key []byte, value []byte) error {
	if isAccountKey(key) {
		accountID := parseAccountID(key)
		balance := parseBalance(value)

		// First try to look up the account
		accounts, err := t.client.LookupAccounts([]types.Uint128{accountID})
		if err != nil {
			return err
		}

		if len(accounts) > 0 {
			// Account exists, but TigerBeetle doesn't support direct updates
			// For now, we'll just verify the balance matches what we expect
			// In a real implementation, you'd use transfers to adjust balances
			if accounts[0].UserData64 != balance {
				// Log a warning that balance doesn't match
				fmt.Printf("Warning: Account %s balance mismatch: %d vs %d\n",
					string(key), accounts[0].UserData64, balance)
			}
			return nil
		} else {
			// Account doesn't exist, create it
			account := types.Account{
				ID:         accountID,
				UserData64: balance, // Store balance in UserData64
				// Set other fields as needed
			}

			result, err := t.client.CreateAccounts([]types.Account{account})
			if err != nil {
				return err
			}

			if len(result) > 0 && result[0].Result != types.AccountOK && result[0].Result != types.AccountExists {
				return fmt.Errorf("failed to create account: %v", result[0].Result)
			}
		}

		return nil
	}

	// return error if key is not an account
	return fmt.Errorf("unsupported key: %s", key)
}

// BeginTx starts a new transaction
func (t *TigerBeetleDB) BeginTx() (Transaction, error) {
	return &TigerBeetleTransaction{
		db:            t,
		pendingWrites: make(map[string][]byte),
	}, nil
}

// Close closes the database
func (t *TigerBeetleDB) Close() error {
	t.client.Close()
	return nil
}

// TigerBeetleTransaction implements the Transaction interface for TigerBeetle
type TigerBeetleTransaction struct {
	db            *TigerBeetleDB
	pendingWrites map[string][]byte
}

// Set stores a key-value pair within a transaction
func (t *TigerBeetleTransaction) Set(key []byte, value []byte) error {
	// Track the write
	t.pendingWrites[string(key)] = value
	return nil
}

// Commit commits the transaction
func (t *TigerBeetleTransaction) Commit() error {
	// First, create any accounts
	for key, value := range t.pendingWrites {
		t.db.Set([]byte(key), value)
	}

	// Clear the pending writes
	t.pendingWrites = make(map[string][]byte)
	return nil
}

// Rollback aborts the transaction
func (t *TigerBeetleTransaction) Rollback() error {
	// Just discard the pending changes
	t.pendingWrites = make(map[string][]byte)
	return nil
}
