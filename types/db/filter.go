package db

import (
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

type BlockFilter struct {
	CreatedAt uint64
	Metadata  primitives.VarData
	From      BlockKey
	Next      BlockKey
	To        BlockKey
}

type TransactionFilter struct {
	CreatedAt uint64
	Metadata  primitives.VarData
	From      TransactionKey
	Next      TransactionKey
	To        TransactionKey
}

type LogFilter struct {
	CreatedAt uint64
	Metadata  primitives.VarData
	From      LogKey
	Next      LogKey
	To        LogKey
	Addresses []primitives.Data20
	Topics    [][]primitives.Data32
}
