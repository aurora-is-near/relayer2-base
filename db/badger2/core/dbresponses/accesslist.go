package dbresponses

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
)

// https://openethereum.github.io/JSONRPC
type AccessListEntry struct {
	Address     dbp.Data20   `json:"address"`
	StorageKeys []dbp.Data32 `json:"storageKeys"`
}
