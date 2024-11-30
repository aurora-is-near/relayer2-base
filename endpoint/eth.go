package endpoint

import (
	"context"
	"errors"

	utils2 "github.com/aurora-is-near/relayer2-base/types/utils"

	"github.com/aurora-is-near/relayer2-base/types"
	"github.com/aurora-is-near/relayer2-base/types/common"
	"github.com/aurora-is-near/relayer2-base/types/engine"
	errs "github.com/aurora-is-near/relayer2-base/types/errors"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/request"
	"github.com/aurora-is-near/relayer2-base/types/response"
	"github.com/aurora-is-near/relayer2-base/utils"
)

type Eth struct {
	*Endpoint
}

func NewEth(endpoint *Endpoint) *Eth {
	return &Eth{endpoint}
}

// Accounts returns empty array
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) Accounts(_ context.Context) (*[]string, error) {
	return utils.Constants.EmptyArray(), nil
}

// Coinbase returns constant 0x0, see relayer.yml to configure coinBase
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) Coinbase(_ context.Context) (*string, error) {
	return &e.Config.EthConfig.ZeroAddress, nil
}

// ProtocolVersion returns constant 0x41, see relayer.yml to configure ProtocolVersion
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) ProtocolVersion(_ context.Context) (*common.Uint256, error) {
	return &e.Config.EthConfig.ProtocolVersion, nil
}

// Hashrate returns constant 0x0, see relayer.yml to configure Hashrate
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) Hashrate(_ context.Context) (*common.Uint256, error) {
	return &e.Config.EthConfig.Hashrate, nil
}

// Mining returns constant false
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) Mining(_ context.Context) (*bool, error) {
	return utils.Constants.Mining(), nil
}

// Syncing returns constant false
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) Syncing(_ context.Context) (*bool, error) {
	return utils.Constants.Syncing(), nil
}

// BlockNumber returns the latest block number from DB if API is enabled by configuration.
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure, returns errors code '-32000' with custom message.
func (e *Eth) BlockNumber(ctx context.Context) (*primitives.HexUint, error) {
	bn, err := e.DbHandler.BlockNumber(ctx)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return bn, nil
}

// GetBlockByHash returns the block from DB, with the given block hash. `hash` is required but `isFull` is an optional
// parameter, if not provided default is false. If isFull is true all transactions in the block with all details
// otherwise returns only the hashes of the transactions are returned.
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On missing or invalid param returns errors code '-32602' with custom message.
// 	If hash not found (KeyNotFoundError) returns nil
// 	On DB failure or hash not found, returns errors code '-32000' with custom message.
func (e *Eth) GetBlockByHash(ctx context.Context, hash common.H256, isFull *bool) (*response.Block, error) {
	if isFull == nil {
		isFull = utils.Constants.Full()
	}
	block, err := e.DbHandler.GetBlockByHash(ctx, hash, *isFull)
	if err != nil {
		_, ok := err.(*errs.KeyNotFoundError)
		if !ok {
			return nil, err
		}
		return nil, nil
	}
	return block, nil
}

// GetBlockByNumber returns the block from DB, with the given block number. `number` is required but `isFull` is an
// optional parameter, if not provided default is false. If isFull is true all transactions in the block with all details
// otherwise returns only the hashes of the transactions are returned.
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or number not found, returns errors code '-32000' with custom message.
//	On missing or invalid param returns errors code '-32602' with custom message.
func (e *Eth) GetBlockByNumber(ctx context.Context, number common.BN64, isFull *bool) (*response.Block, error) {
	if isFull == nil {
		isFull = utils.Constants.Full()
	}
	block, err := e.DbHandler.GetBlockByNumber(ctx, number, *isFull)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return block, nil
}

// GetBlockTransactionCountByHash returns the number of transactions withing the given block hash.
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
//	On DB failure or hash not found, returns errors code '-32000' with custom message.
//	On missing or invalid param returns errors code '-32602' with custom message.
func (e *Eth) GetBlockTransactionCountByHash(ctx context.Context, hash common.H256) (*primitives.HexUint, error) {
	cnt, err := e.DbHandler.GetBlockTransactionCountByHash(ctx, hash)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return cnt, nil
}

// GetBlockTransactionCountByNumber returns the number of transactions within the given block number. `number` parameter
// is optional and the latest block is used if `number` parameter is not provided
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or number not found, returns errors code '-32000' with custom message.
// 	On missing or invalid param returns errors code '-32602' with custom message.
func (e *Eth) GetBlockTransactionCountByNumber(ctx context.Context, number *common.BN64) (*primitives.HexUint, error) {
	bn := common.LatestBlockNumber
	if number != nil {
		bn = *number
	}
	cnt, err := e.DbHandler.GetBlockTransactionCountByNumber(ctx, bn)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return cnt, nil
}

// GetTransactionByHash returns the transaction information of the given transaction hash.
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or hash not found, returns errors code '-32000' with custom message.
// 	On missing or invalid param returns errors code '-32602' with custom message.
func (e *Eth) GetTransactionByHash(ctx context.Context, hash common.H256) (*response.Transaction, error) {
	tx, err := e.DbHandler.GetTransactionByHash(ctx, hash)
	if err != nil {
		_, ok := err.(*errs.KeyNotFoundError)
		if !ok {
			return nil, err
		}
		return nil, nil
	}
	return tx, nil
}

// GetTransactionByBlockHashAndIndex returns the transaction information of the given block hash and transaction index
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or hash not found, returns errors code '-32000' with custom message.
func (e *Eth) GetTransactionByBlockHashAndIndex(ctx context.Context, hash common.H256, index common.Uint64) (*response.Transaction, error) {
	tx, err := e.DbHandler.GetTransactionByBlockHashAndIndex(ctx, hash, index)
	if err != nil {
		_, ok := err.(*errs.KeyNotFoundError)
		if !ok {
			return nil, err
		}
		return nil, nil
	}
	return tx, nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction information of the given block number and transaction index.
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or number not found, returns errors code '-32000' with custom message.
func (e *Eth) GetTransactionByBlockNumberAndIndex(ctx context.Context, number common.BN64, index common.Uint64) (*response.Transaction, error) {
	tx, err := e.DbHandler.GetTransactionByBlockNumberAndIndex(ctx, number, index)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return tx, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash.
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On missing or invalid param returns error code '-32602' with custom message.
// 	If hash not found (KeyNotFoundError) returns nil
// 	On DB failure or other internal errors, returns errors code '-32000' with custom message.
func (e *Eth) GetTransactionReceipt(ctx context.Context, hash common.H256) (*response.TransactionReceipt, error) {
	txsReceipt, err := e.DbHandler.GetTransactionReceipt(ctx, hash)
	if err != nil {
		_, ok := err.(*errs.KeyNotFoundError)
		if !ok {
			return nil, err
		}
		return nil, nil
	}
	return txsReceipt, nil
}

// NewFilter creates a new filter based on the filter options and returns newly created filter ID on success.
//
// FilterOptions object is mandatory but all keys are optional
// 	- fromBlock: QUANTITY|TAG - Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
// 	- toBlock: QUANTITY|TAG - Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
// 	- address: DATA|Array - Contract address or a list of addresses from which logs should originate.
//	- topics: Array of DATA - Topics are order-dependent. Each topic can also be an array of DATA with "or" options.
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On filter option parsing failure, returns errors code '32602' with custom message.
//	On DB failure, returns errors code '-32000' with custom message.
func (e *Eth) NewFilter(ctx context.Context, rawFilter request.Filter) (*common.Uint256, error) {
	f, err := e.parseRequestFilter(ctx, &rawFilter)
	if err != nil {
		return nil, &errs.InvalidParamsError{Message: err.Error()}
	}
	dbf := f.ToLogFilter()
	fid := common.RandomUint256()
	err = e.DbHandler.StoreLogFilter(ctx, fid.Data32(), dbf)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return &fid, nil
}

// NewBlockFilter creates a filter and returns newly created filter ID on success. To check if the state has changed,
// call "eth_getFilterChanges".
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure, returns errors code '-32000' with custom message.
func (e *Eth) NewBlockFilter(ctx context.Context) (*common.Uint256, error) {
	bn, err := e.DbHandler.BlockNumber(ctx)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	fromBlock := uint64(*bn)
	fid := common.RandomUint256()
	nFilter := &(types.Filter{FromBlock: &fromBlock})
	dbf := nFilter.ToBlockFilter()
	err = e.DbHandler.StoreBlockFilter(ctx, fid.Data32(), dbf)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return &fid, nil
}

// NewPendingTransactionFilter returns empty array
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) NewPendingTransactionFilter(_ context.Context) (*string, error) {
	return utils.Constants.ZeroStrUint128(), nil
}

// UninstallFilter deletes a filter with given filter id and returns true on success. Additionally, filters timeout when
// they aren't requested with "eth_getFilterChanges" for a period of time.
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
//	On DB failure or filterId not found, returns errors code '-32000' with custom message.
func (e *Eth) UninstallFilter(ctx context.Context, filterId common.Uint256) (*bool, error) {
	var err error
	resp := true
	err = e.DbHandler.DeleteFilter(ctx, filterId.Data32())
	if err != nil {
		resp = false
	}
	return &resp, err
}

// GetFilterChanges polls method for a filter, on success returns an array of logs which occurred since last poll.
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
//	On failure, returns errors code '-32000' with custom message.
func (e *Eth) GetFilterChanges(ctx context.Context, filterId common.Uint256) (*[]interface{}, error) {
	fid := filterId.Data32()
	filter, err := e.DbHandler.GetFilter(ctx, fid)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	resp, err := e.DbHandler.GetFilterChanges(ctx, filter)
	err = e.DbHandler.StoreFilter(ctx, fid, filter)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return resp, err
}

// GetLogs returns an array of log objects for the given filter
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
//	On filter option parsing failure, returns errors code '32602' with custom message.
// 	On DB failure, returns errors code '-32000' with custom message.
func (e *Eth) GetLogs(ctx context.Context, rawFilter request.Filter) (*[]*response.Log, error) {
	filter, err := e.parseRequestFilter(ctx, &rawFilter)
	if err != nil {
		return nil, &errs.InvalidParamsError{Message: err.Error()}
	}
	dbf := filter.ToLogFilter()
	logResponses, err := e.DbHandler.GetLogs(ctx, dbf)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return &logResponses, nil
}

// GetFilterLogs returns an array of all logs matching filter with given id.
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
//	On DB failure, returns errors code '-32000' with custom message.
func (e *Eth) GetFilterLogs(ctx context.Context, filterId common.Uint256) (*[]*response.Log, error) {
	fid := filterId.Data32()
	filter, err := e.DbHandler.GetLogFilter(ctx, fid)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	logs, err := e.DbHandler.GetLogs(ctx, filter)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	err = e.DbHandler.StoreLogFilter(ctx, fid, filter)
	if err != nil {
		return nil, &errs.GenericError{Err: err}
	}
	return &logs, nil
}

// GetUncleCountByBlockHash returns zero
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On missing or invalid param returns error code '-32602' with custom message.
// 	If hash not found (KeyNotFoundError) returns nil
// 	On DB failure or other internal errors, returns errors code '-32000' with custom message.
func (e *Eth) GetUncleCountByBlockHash(ctx context.Context, hash common.H256) (*common.Uint256, error) {
	block, err := e.DbHandler.GetBlockByHash(ctx, hash, false)
	if block == nil || err != nil {
		_, ok := err.(*errs.KeyNotFoundError)
		if !ok {
			return nil, err
		}
		return nil, nil
	}
	return utils.Constants.ZeroUint256(), nil
}

// GetUncleCountByBlockNumber returns zero
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On missing or invalid param returns error code '-32602' with custom message.
// 	If block number not found (KeyNotFoundError) returns nil
// 	On DB failure or other internal errors, returns errors code '-32000' with custom message.
func (e *Eth) GetUncleCountByBlockNumber(ctx context.Context, number *common.BN64) (*common.Uint256, error) {
	if number != nil {
		block, err := e.DbHandler.GetBlockByNumber(ctx, *number, false)
		if block == nil || err != nil {
			_, ok := err.(*errs.KeyNotFoundError)
			if !ok {
				return nil, err
			}
			return nil, nil
		}
	}
	return utils.Constants.ZeroUint256(), nil
}

// GetUncleByBlockHashAndIndex returns null/nil
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetUncleByBlockHashAndIndex(_ context.Context, _ *common.H256, _ *common.Uint64) (*string, error) {
	return nil, nil
}

// GetUncleByBlockNumberAndIndex returns null/nil
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetUncleByBlockNumberAndIndex(_ context.Context, _ *common.BN64, _ *common.Uint64) (*string, error) {
	return nil, nil
}

// GetCompilers returns empty array
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) GetCompilers(_ context.Context) (*[]string, error) {
	return utils.Constants.EmptyArray(), nil
}

// PendingTransactions returns empty array
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) PendingTransactions(_ context.Context) (*[]string, error) {
	return utils.Constants.EmptyArray(), nil
}

// EstimateGas returns constant gas estimation provided in configuration file.
// The endpoint should be proxied to Mainnet to get an estimate of how much gas is necessary to allow the transaction to complete.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) EstimateGas(_ context.Context, txs engine.TransactionForCall, bNumOrHash *common.BlockNumberOrHash) (*common.Uint256, error) {
	return &e.Config.EthConfig.GasEstimate, nil
}

// GasPrice returns constant gas price provided in the configuration file.
// The endpoint should be proxied to Mainnet to get the value used in the Aurora infrastructure.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) GasPrice(_ context.Context) (*common.Uint256, error) {
	return &e.Config.EthConfig.GasPrice, nil
}

func (e *Eth) parseRequestFilter(ctx context.Context, filter *request.Filter) (*types.Filter, error) {

	f := &types.Filter{}
	if filter.BlockHash != nil {
		bn, err := e.DbHandler.BlockHashToNumber(ctx, *filter.BlockHash)
		if err != nil {
			// Return empty array if Block Hash not found or not converted to valid Block Number
			f.FromBlock, f.ToBlock = nil, nil
			return f, nil
		}
		f.FromBlock, f.ToBlock = bn, bn
	} else {
		if filter.ToBlock != nil {
			f.ToBlock = filter.ToBlock.Uint64()
			// Update block "0x0" (earliest)as "0x1" to prevent zero value and zero index confusion
			if f.ToBlock != nil && *f.ToBlock == 0 {
				*f.ToBlock += 1
			}

		}
		if filter.FromBlock != nil {
			f.FromBlock = filter.FromBlock.Uint64()
			// Update block "0x0" (earliest)as "0x1" to prevent zero value and zero index confusion
			if f.FromBlock != nil && *f.FromBlock == 0 {
				*f.FromBlock += 1
			}
		} else {
			bn, err := e.DbHandler.BlockNumber(ctx)
			if err == nil && bn != nil {
				bn64, err := utils2.HexStringToUint64(bn.Hex())
				if err == nil {
					f.FromBlock = &bn64
				}
			}
		}

		// toBlock specified while fromBlock not => possible from > to case since fromBlock is latest if not specified
		if f.FromBlock == nil && f.ToBlock != nil {
			return nil, errors.New("fromBlock cannot be latest while toBlock is specified")
		}

		// if both specified and from > to then do not save filter at first place
		if f.FromBlock != nil && f.ToBlock != nil {
			if *f.FromBlock > *f.ToBlock {
				return nil, errors.New("fromBlock cannot be greater than toBlock")
			}
		}
	}

	f.Addresses = make([]primitives.Data20, 0)
	if filter.Addresses != nil {
		seen := make(map[string]bool)
		for _, a := range filter.Addresses {
			if seen[a.Hex()] {
				continue
			}
			f.Addresses = append(f.Addresses, primitives.Data20FromBytes(a.Bytes()))
			seen[a.Hex()] = true
		}
	}

	f.Topics = filter.Topics

	return f, nil
}
