//go:build !tigerbeetle

package db

import (
	"fmt"
)

// NewTigerBeetleDBFromMain is a stub function for non-TigerBeetle builds
func NewTigerBeetleDBFromMain(addresses string) (DB, error) {
	return nil, fmt.Errorf("TigerBeetle support requires building with -tags tigerbeetle")
}
