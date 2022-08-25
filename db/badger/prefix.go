package badger

import (
	"aurora-relayer-go-common/utils"
	"bytes"
	"fmt"
)

type TablePrefix string

const (
	prefixCurrentBlockId                   TablePrefix = "/current-block-id"
	prefixBlockByHash                      TablePrefix = "/block/hash/"
	prefixBlockByNumber                    TablePrefix = "/block/number/"
	prefixTransactionCountByBlockHash      TablePrefix = "/transaction-count/hash/"
	prefixTransactionCountByBlockNumber    TablePrefix = "/transaction-count/number/"
	prefixTransactionByBlockHashAndIndex   TablePrefix = "/transaction/block-hash/"
	prefixTransactionByBlockNumberAndIndex TablePrefix = "/transaction/block-num/"
	prefixTransactionByHash                TablePrefix = "/transaction/hash/"
	prefixLogTable                         TablePrefix = "/logs/"
	prefixLogIndexTable                    TablePrefix = "/logs-index/"
)

type Keyable interface {
	KeyBytes() []byte
}

func (t TablePrefix) AppendBytes(args ...[]byte) []byte {
	k := t.Bytes()
	for _, arg := range args {
		k = append(k, arg...)
	}
	return k
}

func (t TablePrefix) Key(args ...Keyable) []byte {
	k := t.Bytes()
	if len(args) == 0 {
		return k
	} else if len(args) == 1 {
		return append(k, args[0].KeyBytes()...)
	} else {
		for _, arg := range args {
			k = append(k, arg.KeyBytes()...)
		}
		return k
	}
}

func (t TablePrefix) Bytes() []byte {
	return []byte(t)
}

func getLogInsertKey(blockNum utils.Uint256, txIdx, logIdx Keyable) []byte {
	buf := new(bytes.Buffer)
	_, _ = buf.Write(blockNum.KeyBytes())
	_, _ = buf.Write(txIdx.KeyBytes())
	_, _ = buf.Write(logIdx.KeyBytes())
	return buf.Bytes()
}

func filterByIdKey(id utils.Uint256) []byte {
	return []byte(fmt.Sprintf("/filter-by-id/%s", id))
}
