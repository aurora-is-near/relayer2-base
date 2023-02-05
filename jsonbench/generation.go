package jsonbench

import (
	"math/rand"
	"relayer2-base/types/primitives"
	"relayer2-base/types/response"
)

const numBlocks = 10
const numTxs = 10
const numTxReceipts = 10
const numLogs = 10
const numHashArrays = 10

type MyRandom struct {
	r *rand.Rand
}

func NewRandom() MyRandom {
	return MyRandom{r: rand.New(rand.NewSource(528491))}
}

func (mr MyRandom) Get(n int) []byte {
	data := make([]byte, n)
	mr.r.Read(data)
	return data
}

type Payload struct {
	Blocks []*response.Block
}

func generateAccessListEntry(r MyRandom) *response.AccessListEntry {
	ace := &response.AccessListEntry{
		Address: primitives.Data20FromBytes(r.Get(20)),
	}
	n := r.r.Intn(10)
	for i := 0; i < n; i++ {
		ace.StorageKeys = append(ace.StorageKeys, primitives.Data32FromBytes(r.Get(32)))
	}
	return ace
}

func generateTx(r MyRandom) *response.Transaction {
	tx := &response.Transaction{
		BlockHash:        primitives.Data32FromBytes(r.Get(32)),
		BlockNumber:      primitives.HexUint(r.r.Uint64()),
		From:             primitives.Data20FromBytes(r.Get(20)),
		Gas:              primitives.HexUint(r.r.Uint64()),
		GasPrice:         primitives.QuantityFromBytes(r.Get(32)),
		Hash:             primitives.Data32FromBytes(r.Get(32)),
		Input:            primitives.VarDataFromBytes(r.Get(r.r.Intn(100))),
		Nonce:            primitives.QuantityFromBytes(r.Get(32)),
		V:                primitives.HexUint(r.r.Uint64()),
		R:                primitives.QuantityFromBytes(r.Get(32)),
		S:                primitives.QuantityFromBytes(r.Get(32)),
		TransactionIndex: primitives.HexUint(r.r.Uint64()),
		Value:            primitives.QuantityFromBytes(r.Get(32)),
	}
	if r.r.Intn(10) == 0 {
		al := []*response.AccessListEntry{}
		n := r.r.Intn(10)
		for i := 0; i < n; i++ {
			al = append(al, generateAccessListEntry(r))
		}
		tx.AccessList = &al
	}
	if r.r.Intn(2) == 0 {
		chainID := primitives.HexUint(r.r.Uint64())
		tx.ChainID = &chainID
	}
	if r.r.Intn(5) == 0 {
		mpfpg := primitives.QuantityFromBytes(r.Get(32))
		mfpg := primitives.QuantityFromBytes(r.Get(32))
		tx.MaxPriorityFeePerGas = &mpfpg
		tx.MaxFeePerGas = &mfpg
	}
	if r.r.Intn(10) > 0 {
		to := primitives.Data20FromBytes(r.Get(20))
		tx.To = &to
	}
	return tx
}

func generateBlock(r MyRandom) *response.Block {
	block := &response.Block{
		Number:           primitives.HexUint(r.r.Uint64()),
		Hash:             primitives.Data32FromBytes(r.Get(32)),
		ParentHash:       primitives.Data32FromBytes(r.Get(32)),
		Nonce:            primitives.Data8FromBytes(r.Get(8)),
		Sha3Uncles:       primitives.Data32FromBytes(r.Get(32)),
		LogsBloom:        primitives.Data256FromBytes(r.Get(256)),
		TransactionsRoot: primitives.Data32FromBytes(r.Get(32)),
		StateRoot:        primitives.Data32FromBytes(r.Get(32)),
		ReceiptsRoot:     primitives.Data32FromBytes(r.Get(32)),
		Miner:            primitives.Data20FromBytes(r.Get(20)),
		MixHash:          primitives.Data32FromBytes(r.Get(32)),
		Difficulty:       primitives.HexUint(r.r.Uint64()),
		TotalDifficulty:  primitives.HexUint(r.r.Uint64()),
		ExtraData:        primitives.VarDataFromBytes(r.Get(r.r.Intn(100))),
		Size:             primitives.HexUint(r.r.Uint64()),
		GasLimit:         primitives.QuantityFromBytes(r.Get(32)),
		GasUsed:          primitives.QuantityFromBytes(r.Get(32)),
		Timestamp:        primitives.HexUint(r.r.Uint64()),
	}
	n := r.r.Intn(30)
	if r.r.Intn(2) == 0 {
		for i := 0; i < n; i++ {
			block.Transactions = append(block.Transactions, primitives.Data32FromBytes(r.Get(32)))
		}
	} else {
		for i := 0; i < n; i++ {
			block.Transactions = append(block.Transactions, generateTx(r))
		}
	}
	n = r.r.Intn(9)
	for i := 0; i < n; i++ {
		block.Uncles = append(block.Uncles, primitives.Data32FromBytes(r.Get(32)))
	}
	return block
}

func generateLog(r MyRandom) *response.Log {
	log := &response.Log{
		Removed:          r.r.Intn(2) == 0,
		LogIndex:         primitives.HexUint(r.r.Uint64()),
		TransactionIndex: primitives.HexUint(r.r.Uint64()),
		TransactionHash:  primitives.Data32FromBytes(r.Get(32)),
		BlockHash:        primitives.Data32FromBytes(r.Get(32)),
		BlockNumber:      primitives.HexUint(r.r.Uint64()),
		Address:          primitives.Data20FromBytes(r.Get(20)),
		Data:             primitives.VarDataFromBytes(r.Get(r.r.Intn(100))),
	}
	n := r.r.Intn(5)
	for i := 0; i < n; i++ {
		log.Topics = append(log.Topics, primitives.Data32FromBytes(r.Get(32)))
	}
	return log
}

func generateTxReceipt(r MyRandom) *response.TransactionReceipt {
	txReceipt := &response.TransactionReceipt{
		BlockHash:           primitives.Data32FromBytes(r.Get(32)),
		BlockNumber:         primitives.HexUint(r.r.Uint64()),
		CumulativeGasUsed:   primitives.QuantityFromBytes(r.Get(32)),
		From:                primitives.Data20FromBytes(r.Get(20)),
		GasUsed:             primitives.HexUint(r.r.Uint64()),
		LogsBloom:           primitives.Data256FromBytes(r.Get(256)),
		NearReceiptHash:     primitives.Data32FromBytes(r.Get(32)),
		NearTransactionHash: primitives.Data32FromBytes(r.Get(32)),
		Status:              primitives.HexUint(r.r.Uint64()),
		TransactionHash:     primitives.Data32FromBytes(r.Get(32)),
		TransactionIndex:    primitives.HexUint(r.r.Uint64()),
	}
	caOrTo := primitives.Data20FromBytes(r.Get(20))
	if r.r.Intn(10) == 0 {
		txReceipt.ContractAddress = &caOrTo
	} else {
		txReceipt.To = &caOrTo
	}
	n := r.r.Intn(10)
	for i := 0; i < n; i++ {
		txReceipt.Logs = append(txReceipt.Logs, generateLog(r))
	}
	return txReceipt
}

func generateHashArray(r MyRandom) []primitives.Data32 {
	hashes := []primitives.Data32{}
	n := r.r.Intn(50)
	for i := 0; i < n; i++ {
		hashes = append(hashes, primitives.Data32FromBytes(r.Get(32)))
	}
	return hashes
}

func generatePayloads(r MyRandom) []any {
	res := []any{}

	for i := 0; i < numBlocks; i++ {
		res = append(res, generateBlock(r))
	}

	for i := 0; i < numTxs; i++ {
		res = append(res, generateTx(r))
	}

	for i := 0; i < numTxReceipts; i++ {
		res = append(res, generateTxReceipt(r))
	}

	for i := 0; i < numLogs; i++ {
		res = append(res, generateLog(r))
	}

	for i := 0; i < numHashArrays; i++ {
		res = append(res, generateHashArray(r))
	}

	r.r.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})

	return res
}
