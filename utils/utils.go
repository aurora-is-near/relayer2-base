package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/tinypack"
	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/indexer"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"golang.org/x/crypto/sha3"
)

const (
	AccountId = "aurora"
)

func IndexerBlockToDbBlock(block *indexer.Block) *dbt.Block {
	b := dbt.Block{
		ParentHash:       block.ParentHash,
		Miner:            block.Miner,
		Timestamp:        uint64(block.Timestamp),
		GasLimit:         block.GasLimit,
		GasUsed:          block.GasUsed,
		LogsBloom:        block.LogsBloom,
		TransactionsRoot: block.TransactionsRoot,
		StateRoot:        block.StateRoot,
		ReceiptsRoot:     block.ReceiptsRoot,
		Size:             uint64(block.Size),
	}
	return &b
}

func IndexerTxnToDbTxn(txn *indexer.Transaction, cumulativeGas primitives.Quantity) *dbt.Transaction {
	toOrContract := tinypack.CreateNullable[primitives.Data20](nil)
	isContractDeployment := false
	if txn.ContractAddress != nil {
		if txn.To != nil {
			log.Log().Warn().Msgf("both contract address and to address is set for txn: [%v], to: [%s], contract: [%s]",
				txn.Hash, txn.To.Hex(), txn.ContractAddress.Hex())
		}
		isContractDeployment = true
		toOrContract = tinypack.CreateNullable[primitives.Data20](txn.ContractAddress)
	} else if txn.To != nil {
		toOrContract = tinypack.CreateNullable[primitives.Data20](txn.To)
	} else {
		log.Log().Warn().Msgf("both contract address and to address is null for txn: [%v]", txn.Hash.Hex())
	}

	var accessListEntries []dbt.AccessListEntry
	for _, al := range txn.AccessList {
		var storageKeys []primitives.Data32
		for _, sk := range al.StorageKeys {
			storageKeys = append(storageKeys, sk)
		}
		sk := tinypack.CreateList[primitives.VarLen, primitives.Data32](storageKeys...)
		accessListEntry := dbt.AccessListEntry{
			Address:     al.Address,
			StorageKeys: tinypack.VarList[primitives.Data32]{sk},
		}
		accessListEntries = append(accessListEntries, accessListEntry)
	}
	ake := tinypack.CreateList[primitives.VarLen, dbt.AccessListEntry](accessListEntries...)
	nearHash := tinypack.CreateNullable[primitives.Data32](nil)
	if txn.NearTransaction.Hash != nil {
		nh := primitives.Data32(*txn.NearTransaction.Hash)
		nearHash = tinypack.CreateNullable[primitives.Data32](&nh)
	}

	t := dbt.Transaction{
		Type:                 txn.TxType,
		From:                 txn.From,
		IsContractDeployment: isContractDeployment,
		ToOrContract:         toOrContract,
		Nonce:                txn.Nonce,
		GasPrice:             txn.GasPrice,
		GasLimit:             txn.GasLimit,
		GasUsed:              txn.GasUsed,
		CumulativeGasUsed:    cumulativeGas,
		Value:                txn.Value,
		Input:                primitives.VarData(txn.Input),
		NearHash:             nearHash,
		NearReceiptHash:      primitives.Data32(txn.NearTransaction.ReceiptHash),
		Status:               txn.Status,
		V:                    txn.V,
		R:                    txn.R,
		S:                    txn.S,
		LogsBloom:            txn.LogsBloom,
		AccessList:           tinypack.VarList[dbt.AccessListEntry]{ake},
		MaxPriorityFeePerGas: txn.MaxPriorityFeePerGas,
		MaxFeePerGas:         txn.MaxFeePerGas,
	}
	return &t
}

func IndexerLogToDbLog(log *indexer.Log) *dbt.Log {
	var topics []primitives.Data32
	for _, t := range log.Topics {
		topics = append(topics, primitives.Data32(t))
	}
	t := tinypack.CreateList[primitives.VarLen, primitives.Data32](topics...)

	l := dbt.Log{
		Address: log.Address,
		Data:    primitives.VarData(log.Data),
		Topics:  tinypack.VarList[primitives.Data32]{t},
	}
	return &l
}

func ComputeBlockHash(bHeight, chainId uint64) []byte {
	bufEmpty25 := make([]byte, 25)

	bufCId := make([]byte, 8)
	binary.BigEndian.PutUint64(bufCId, chainId)

	bufAId := []byte(AccountId)

	bufBH := make([]byte, 8)
	binary.BigEndian.PutUint64(bufBH, bHeight)

	bufHash := append(bufEmpty25, bufCId...)
	bufHash = append(bufHash, bufAId...)
	bufHash = append(bufHash, bufBH...)
	hash := sha256.Sum256(bufHash)

	return hash[:]
}

// ParseEVMRevertReason resolves the abi-encoded revert reason
func ParseEVMRevertReason(data []byte) (string, error) {
	if len(data) < 4 {
		return "", errors.New("invalid data for unpacking")
	}
	// The first 4 bytes (08c379a0) are the function selector for error signature
	errorSig := []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
	if !bytes.Equal(data[:4], errorSig) {
		return "txs result not Error(string)", errors.New("txs result not of type Error(string)")
	}

	reason, err := parseEVMRevertReason(data[4:])
	if err != nil {
		return "invalid txs result", err
	}
	return reason, err
}

// CalculateKeccak256 calculates and returns the Keccak256 hash of the input data
func CalculateKeccak256(input []byte) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(input)
	keccak256 := hash.Sum(nil)
	return fmt.Sprintf("0x%s", hex.EncodeToString(keccak256))
}

// parseEVMRevertReason decodes input slice according to ABI, specification, interprets a 32 byte slice as an offset and
// then determines which indices to look for string start and string length
func parseEVMRevertReason(input []byte) (reason string, err error) {
	bigOffsetEnd := new(big.Int).SetBytes(input[:32])
	bigOffsetEnd.Add(bigOffsetEnd, big.NewInt(32))
	inputLength := big.NewInt(int64(len(input)))

	if bigOffsetEnd.Cmp(inputLength) > 0 {
		return "", fmt.Errorf("offset %v would go over slice boundary (len=%v)", bigOffsetEnd, inputLength)
	}

	if bigOffsetEnd.BitLen() > 63 {
		return "", fmt.Errorf("offset larger than int64: %v", bigOffsetEnd)
	}

	offsetEnd := int(bigOffsetEnd.Uint64())
	bigLength := new(big.Int).SetBytes(input[offsetEnd-32 : offsetEnd])

	totalLength := new(big.Int).Add(bigOffsetEnd, bigLength)
	if totalLength.BitLen() > 63 {
		return "", fmt.Errorf("length larger than int64: %v", totalLength)
	}

	if totalLength.Cmp(inputLength) > 0 {
		return "", fmt.Errorf("length insufficient %v require %v", inputLength, totalLength)
	}

	start := int(bigOffsetEnd.Uint64())
	length := int(bigLength.Uint64())

	return string(input[start : start+length]), nil
}
