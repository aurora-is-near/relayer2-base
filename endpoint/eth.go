package endpoint

import (
	"aurora-relayer-go-common/utils"
	"context"
	"fmt"
)

var (
	protocolVersion = fmt.Sprintf("%#x", 0x41)
	zeroAddress     = fmt.Sprintf("0x%040x", 0)
)

type Eth struct {
	*Endpoint
}

func NewEth(endpoint *Endpoint) *Eth {
	return &Eth{endpoint}
}

func (e *Eth) Accounts(_ context.Context) ([]string, error) {
	if err := e.IsEndpointAllowed("eth_accounts"); err != nil {
		return nil, err
	}
	// TODO implement
	return nil, nil
}

func (e *Eth) Coinbase(_ context.Context) (string, error) {
	if err := e.IsEndpointAllowed("eth_coinbase"); err != nil {
		return "", err
	}
	return zeroAddress, nil
}

func (e *Eth) ProtocolVersion(_ context.Context) (string, error) {
	if err := e.IsEndpointAllowed("eth_protocolVersion"); err != nil {
		return "", err
	}
	return protocolVersion, nil
}

func (e *Eth) Hashrate(_ context.Context) (string, error) {
	if err := e.IsEndpointAllowed("eth_hashrate"); err != nil {
		return "", err
	}
	return utils.IntToHex(0), nil
}

// BlockNumber returns the latest block number from DB if API is enabled by configuration.
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure, returns error code '-32000' with custom message.
func (e *Eth) BlockNumber(_ context.Context) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_blockNumber"); err != nil {
		return nil, err
	}
	bn, err := (*e.DbHandler).BlockNumber()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	number := utils.Uint256(fmt.Sprint(bn))
	return &number, nil
}

// GetBlockByHash returns the block from DB, with the given block hash, both hash and isFull are required.
// If isFull is true all transactions in the block with all details otherwise returns only the hashes of the
// transactions are returned.
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure or hash not found, returns error code '-32000' with custom message.
//
// On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetBlockByHash(_ context.Context, hash utils.H256, isFull bool) (*utils.BlockResponse, error) {
	if err := e.IsEndpointAllowed("eth_getBlockByHash"); err != nil {
		return nil, err
	}
	block, err := (*e.DbHandler).GetBlockByHash(hash)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return block.ToResponse(isFull), nil
}

// GetBlockByNumber returns the block from DB, with the given block number, both number and isFull are required.
// If isFull is true all transactions in the block with all details otherwise returns only the hashes of the
// transactions are returned.
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure or number not found, returns error code '-32000' with custom message.
//
// On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetBlockByNumber(_ context.Context, number utils.Uint256, isFull bool) (*utils.BlockResponse, error) {
	if err := e.IsEndpointAllowed("eth_getBlockByNumber"); err != nil {
		return nil, err
	}
	block, err := (*e.DbHandler).GetBlockByNumber(number)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return block.ToResponse(isFull), nil
}

// GetBlockTransactionCountByHash returns the number of transactions withing the given block hash.
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure or hash not found, returns error code '-32000' with custom message.
//
// On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetBlockTransactionCountByHash(_ context.Context, hash utils.H256) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_getBlockTransactionCountByHash"); err != nil {
		return nil, err
	}
	cnt, err := (*e.DbHandler).GetBlockTransactionCountByHash(hash)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	count := utils.Uint256(fmt.Sprint(cnt))
	return &count, nil
}

// GetBlockTransactionCountByNumber returns the number of transactions within the given block number.
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure or number not found, returns error code '-32000' with custom message.
//
// On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetBlockTransactionCountByNumber(_ context.Context, number utils.Uint256) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_getBlockTransactionCountByNumber"); err != nil {
		return nil, err
	}
	cnt, err := (*e.DbHandler).GetBlockTransactionCountByNumber(number)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	count := utils.Uint256(fmt.Sprint(cnt))
	return &count, nil
}

// GetTransactionByHash returns the transaction information of the given transaction hash.
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure or hash not found, returns error code '-32000' with custom message.
//
// On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetTransactionByHash(_ context.Context, hash utils.H256) (*utils.TransactionResponse, error) {
	if err := e.IsEndpointAllowed("eth_getTransactionByHash"); err != nil {
		return nil, err
	}
	tx, err := (*e.DbHandler).GetTransactionByHash(hash)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return tx.ToResponse(), nil
}

// GetTransactionByBlockHashAndIndex returns the transaction information of the given block hash and transaction index
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure or hash not found, returns error code '-32000' with custom message.
func (e *Eth) GetTransactionByBlockHashAndIndex(_ context.Context, hash utils.H256, index utils.Uint256) (*utils.TransactionResponse, error) {
	if err := e.IsEndpointAllowed("eth_getTransactionByBlockHashAndIndex"); err != nil {
		return nil, err
	}
	idx, err := index.ToInt64()
	if err != nil {
		return nil, &utils.InvalidParamsError{Message: fmt.Sprintf("invalid argument: %s", err.Error())}
	}
	tx, err := (*e.DbHandler).GetTransactionByBlockHashAndIndex(hash, idx)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return tx.ToResponse(), nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction information of the given block number and transaction index
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure or number not found, returns error code '-32000' with custom message.
func (e *Eth) GetTransactionByBlockNumberAndIndex(_ context.Context, number, index utils.Uint256) (*utils.TransactionResponse, error) {
	if err := e.IsEndpointAllowed("eth_getTransactionByBlockNumberAndIndex"); err != nil {
		return nil, err
	}
	idx, err := index.ToInt64()
	if err != nil {
		return nil, &utils.InvalidParamsError{Message: fmt.Sprintf("invalid argument: %s", err.Error())}
	}
	tx, err := (*e.DbHandler).GetTransactionByBlockNumberAndIndex(number, idx)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return tx.ToResponse(), nil
}

func (e *Eth) GetLogs(_ context.Context, addr utils.Address, bn utils.Uint256, topic ...[]string) (*utils.Log, error) {
	if err := e.IsEndpointAllowed("eth_getLogs"); err != nil {
		return nil, err
	}
	// TODO implement
	return nil, nil
}
