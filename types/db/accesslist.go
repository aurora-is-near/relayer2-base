package db

import (
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

type AccessListEntry struct {
	Address     primitives.Data20
	StorageKeys []primitives.Data32
}
