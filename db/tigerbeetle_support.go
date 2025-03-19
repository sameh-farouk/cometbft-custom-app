//go:build tigerbeetle

package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

// NewTigerBeetleDBFromMain creates a new TigerBeetleDB instance from a comma-separated list of addresses
func NewTigerBeetleDBFromMain(addresses string) (DB, error) {
	addressList := strings.Split(addresses, ",")
	return NewTigerBeetleDB(addressList)
}

// InitializeTigerBeetleAccounts pre-creates accounts in TigerBeetle
func InitializeTigerBeetleAccounts(db DB, accountIDs []string, initialBalances map[string]uint64) error {
	// Type assertion to get the TigerBeetleDB instance
	tbDB, ok := db.(*TigerBeetleDB)
	if !ok {
		return fmt.Errorf("database is not a TigerBeetleDB instance")
	}

	accounts := make([]types.Account, 0, len(accountIDs))

	for _, id := range accountIDs {
		accountID := parseAccountID([]byte(id))
		balance := initialBalances[id]

		account := types.Account{
			ID:         accountID,
			UserData64: balance,
			// Set other fields as needed
		}

		accounts = append(accounts, account)
	}

	if len(accounts) > 0 {
		result, err := tbDB.client.CreateAccounts(accounts)
		if err != nil {
			return err
		}

		for i, res := range result {
			if res.Result != types.AccountOK && res.Result != types.AccountExists {
				return fmt.Errorf("failed to create account %s: %v", accountIDs[i], res.Result)
			}
		}

		log.Printf("Successfully initialized %d TigerBeetle accounts", len(accounts))
	}

	return nil
}
