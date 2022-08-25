package endpoint

import (
	"aurora-relayer-go-common/utils"
	"context"
)

var (
	protocolVersion = utils.IntToUint256(0x41)
	hashrate        = utils.IntToUint256(0)
	zero            = utils.IntToUint256(0)
)

type Eth struct {
	*Endpoint
}

func NewEth(endpoint *Endpoint) *Eth {
	return &Eth{endpoint}
}

// Accounts returns empty array
//
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) Accounts(_ context.Context) ([]string, error) {
	return []string{}, nil
}

// Coinbase returns constant 0x0
//
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) Coinbase(_ context.Context) (*utils.Uint256, error) {
	return &hashrate, nil
}

// ProtocolVersion returns constant 0x41
//
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) ProtocolVersion(_ context.Context) (*utils.Uint256, error) {
	return &protocolVersion, nil
}

// Hashrate returns constant 0x0
//
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
func (e *Eth) Hashrate(_ context.Context) (string, error) {
	return utils.IntToHex(0), nil
}

// BlockNumber returns the latest block number from DB if API is enabled by configuration.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure, returns error code '-32000' with custom message.
func (e *Eth) BlockNumber(_ context.Context) (*utils.Uint256, error) {
	bn, err := (*e.DbHandler).BlockNumber()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return bn, nil
}

// GetBlockByHash returns the block from DB, with the given block hash, both hash and isFull are required.
// If isFull is true all transactions in the block with all details otherwise returns only the hashes of the
// transactions are returned.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or hash not found, returns error code '-32000' with custom message.
// 	On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetBlockByHash(_ context.Context, hash utils.H256, isFull bool) (*utils.BlockResponse, error) {
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
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or number not found, returns error code '-32000' with custom message.
//	On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetBlockByNumber(_ context.Context, number utils.Uint256, isFull bool) (*utils.BlockResponse, error) {
	block, err := (*e.DbHandler).GetBlockByNumber(number)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return block.ToResponse(isFull), nil
}

// GetBlockTransactionCountByHash returns the number of transactions withing the given block hash.
//
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//	On DB failure or hash not found, returns error code '-32000' with custom message.
//	On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetBlockTransactionCountByHash(_ context.Context, hash utils.H256) (utils.Uint256, error) {
	var count utils.Uint256
	cnt, err := (*e.DbHandler).GetBlockTransactionCountByHash(hash)
	if err != nil {
		return count, &utils.GenericError{Err: err}
	}
	count = utils.IntToUint256(cnt)
	return count, nil
}

// GetBlockTransactionCountByNumber returns the number of transactions within the given block number.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or number not found, returns error code '-32000' with custom message.
// 	On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetBlockTransactionCountByNumber(_ context.Context, number utils.Uint256) (utils.Uint256, error) {
	var count utils.Uint256
	cnt, err := (*e.DbHandler).GetBlockTransactionCountByNumber(number)
	if err != nil {
		return count, &utils.GenericError{Err: err}
	}
	count = utils.IntToUint256(cnt)
	return count, nil
}

// GetTransactionByHash returns the transaction information of the given transaction hash.
//
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or hash not found, returns error code '-32000' with custom message.
// 	On missing or invalid param returns error code '-32602' with custom message.
func (e *Eth) GetTransactionByHash(_ context.Context, hash utils.H256) (*utils.TransactionResponse, error) {
	tx, err := (*e.DbHandler).GetTransactionByHash(hash)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return tx.ToResponse(), nil
}

// GetTransactionByBlockHashAndIndex returns the transaction information of the given block hash and transaction index
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or hash not found, returns error code '-32000' with custom message.
func (e *Eth) GetTransactionByBlockHashAndIndex(_ context.Context, hash utils.H256, index utils.Uint256) (*utils.TransactionResponse, error) {
	tx, err := (*e.DbHandler).GetTransactionByBlockHashAndIndex(hash, index)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return tx.ToResponse(), nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction information of the given block number and transaction index.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or number not found, returns error code '-32000' with custom message.
func (e *Eth) GetTransactionByBlockNumberAndIndex(_ context.Context, number, index utils.Uint256) (*utils.TransactionResponse, error) {
	tx, err := (*e.DbHandler).GetTransactionByBlockNumberAndIndex(number, index)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return tx.ToResponse(), nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or hash not found, returns error code '-32000' with custom message.
func (e *Eth) GetTransactionReceipt(_ context.Context, hash utils.H256) (*utils.TransactionReceiptResponse, error) {
	tx, err := (*e.DbHandler).GetTransactionByHash(hash)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	logs, err := (*e.DbHandler).GetLogsForTransaction(tx)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	resp := tx.ToReceiptResponse()
	resp.Logs = logs
	return resp, nil
}

// GetLogs returns an array of log objects for the given filter
//
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure or number not found, returns error code '-32000' with custom message.
func (e *Eth) GetLogs(_ context.Context, rawFilter utils.FilterOptions) (*[]utils.LogResponse, error) {
	filter, err := e.formatFilterOptions(&rawFilter)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	logResponses, err := (*e.DbHandler).GetLogs(*filter)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return logResponses, nil
}

// NewFilter creates a new filter based on the filter options and returns newly created filter ID on success.
//
// FilterOptions object is mandatory but all keys are optional
// 	- fromBlock: QUANTITY|TAG - Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
// 	- toBlock: QUANTITY|TAG - Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
// 	- address: DATA|Array - Contract address or a list of addresses from which logs should originate.
//	- topics: Array of DATA - Topics are order-dependent. Each topic can also be an array of DATA with "or" options.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On filter option parsing failure, returns error code '32602' with custom message.
//	On DB failure, returns error code '-32000' with custom message.
func (e *Eth) NewFilter(_ context.Context, filterOptions utils.FilterOptions) (*utils.Uint256, error) {
	parsedFilter, err := e.formatFilterOptions(&filterOptions)
	if err != nil {
		return nil, &utils.InvalidParamsError{Message: err.Error()}
	}
	bn, err := (*e.DbHandler).BlockNumber()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	filterId, err := utils.RandomUint256()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	filter := &utils.StoredFilter{
		Type:      "event",
		CreatedBy: "0.0.0.0",
		// BlockHash: parsedFilter.BlockHash,
		FromBlock: parsedFilter.FromBlock,
		ToBlock:   parsedFilter.ToBlock,
		Addresses: parsedFilter.Address,
		Topics:    parsedFilter.Topics,
		PollBlock: *bn,
	}
	(*e.DbHandler).StoreFilter(*filterId, filter)
	if err != nil {
		return nil, err
	}
	return filterId, nil
}

// NewBlockFilter creates a filter and returns newly created filter ID on success. To check if the state has changed,
// call "eth_getFilterChanges".
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On DB failure, returns error code '-32000' with custom message.
func (e *Eth) NewBlockFilter(_ context.Context) (*utils.Uint256, error) {
	bn, err := (*e.DbHandler).BlockNumber()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	filterId, err := utils.RandomUint256()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	filter := &utils.StoredFilter{
		Type:      "block",
		CreatedBy: "0.0.0.0",
		PollBlock: *bn.Add(1),
	}
	err = (*e.DbHandler).StoreFilter(*filterId, filter)
	if err != nil {
		return nil, err
	}
	return filterId, nil
}

// UninstallFilter deletes a filter with given filter id and returns true on success. Additionally, filters timeout when
// they aren't requested with "eth_getFilterChanges" for a period of time.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On failure returns false
func (e *Eth) UninstallFilter(_ context.Context, filterId utils.Uint256) (bool, error) {
	err := (*e.DbHandler).DeleteFilter(filterId)
	if err != nil {
		e.Logger.Err(err).Msgf("failed to uninstall filter [%d]", filterId)
		return false, nil
	}
	return true, nil
}

// GetFilterChanges polls method for a filter, on success returns an array of logs which occurred since last poll.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//	On failure, returns error code '-32000' with custom message.
func (e *Eth) GetFilterChanges(_ context.Context, filterId utils.Uint256) (*[]interface{}, error) {

	storedFilter, err := (*e.DbHandler).GetFilter(filterId)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}

	filterChanges := make([]interface{}, 0)

	switch storedFilter.Type {
	case "event":
		logFilter := utils.LogFilter{
			FromBlock: storedFilter.FromBlock,
			ToBlock:   storedFilter.ToBlock,
			Address:   storedFilter.Addresses,
			Topics:    storedFilter.Topics,
		}
		logs, err := (*e.DbHandler).GetLogs(logFilter)
		if err != nil {
			return nil, &utils.GenericError{Err: err}
		}
		for _, l := range *logs {
			filterChanges = append(filterChanges, l)
		}

	case "block":
		blocks, err := (*e.DbHandler).GetBlockHashesSinceNumber(storedFilter.PollBlock)
		if err != nil {
			return nil, &utils.GenericError{Err: err}
		}
		for _, b := range blocks {
			filterChanges = append(filterChanges, b)
		}
	case "transaction":
	default:
	}

	bn, err := (*e.DbHandler).BlockNumber()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	storedFilter.PollBlock = *bn.Add(1)
	(*e.DbHandler).StoreFilter(filterId, storedFilter)

	return &filterChanges, nil
}

// GetFilterLogs returns an array of all logs matching filter with given id.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//	On failure, returns error code '-32000' with custom message.
func (e *Eth) GetFilterLogs(_ context.Context, filterId utils.Uint256) (*[]interface{}, error) {

	storedFilter, err := (*e.DbHandler).GetFilter(filterId)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}

	var res = make([]interface{}, 0)

	switch storedFilter.Type {
	case "event":
		logFilter := utils.LogFilter{
			FromBlock: storedFilter.FromBlock,
			ToBlock:   storedFilter.ToBlock,
			Address:   storedFilter.Addresses,
			Topics:    storedFilter.Topics,
		}
		logs, err := (*e.DbHandler).GetLogs(logFilter)
		if err != nil {
			return nil, &utils.GenericError{Err: err}
		}
		for _, l := range *logs {
			res = append(res, l)
		}

	case "block":
		blocks, err := (*e.DbHandler).GetBlockHashesSinceNumber(storedFilter.PollBlock)
		if err != nil {
			return nil, &utils.GenericError{Err: err}
		}
		for _, b := range blocks {
			res = append(res, b)
		}

	case "transaction":
	default:
	}

	curr, err := (*e.DbHandler).BlockNumber()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	storedFilter.PollBlock = *curr.Add(1)
	err = (*e.DbHandler).StoreFilter(filterId, storedFilter)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}

	return &res, nil
}

// GetUncleCountByBlockHash
//
// 	TODO: implement
func (e *Eth) GetUncleCountByBlockHash(_ context.Context, hash utils.H256) (*utils.Uint256, error) {
	_, err := (*e.DbHandler).GetBlockByHash(hash)
	if err != nil {
		return nil, nil
	}
	return &zero, nil
}

// GetUncleCountByBlockNumber
//
//	TODO: implement
func (e *Eth) GetUncleCountByBlockNumber(_ context.Context, number utils.Uint256) (*utils.Uint256, error) {
	_, err := (*e.DbHandler).GetBlockByNumber(number)
	if err != nil {
		return nil, nil
	}
	return &zero, nil
}

func (e *Eth) formatFilterOptions(filterOptions *utils.FilterOptions) (*utils.LogFilter, error) {
	result := &utils.LogFilter{}
	if filterOptions.BlockHash != nil {
		blockNum, err := (*e.DbHandler).BlockHashToNumber(*filterOptions.BlockHash)
		if err != nil {
			return nil, err
		}
		result.FromBlock, result.ToBlock = blockNum, blockNum
	} else {
		var err error
		result.FromBlock, result.ToBlock, err = e.parseFromAndTo(filterOptions)
		if err != nil {
			return nil, err
		}
	}
	if filterOptions.Address != nil {
		result.Address = make([][]byte, 0)
		seen := make(map[utils.Address]bool)
		for _, a := range filterOptions.Address {
			if seen[a] {
				continue
			}
			result.Address = append(result.Address, a.Bytes())
			seen[a] = true
		}
	}
	result.Topics = filterOptions.Topics
	return result, nil
}

func (e *Eth) parseFromAndTo(filter *utils.FilterOptions) (from, to *utils.Uint256, err error) {
	from, err = parseBlock(filter.FromBlock)
	if err != nil {
		return nil, nil, err
	}
	to, err = parseBlock(filter.ToBlock)
	if err != nil {
		return nil, nil, err
	}

	if from != nil && to != nil {
		return
	}

	latest, err := (*e.DbHandler).BlockNumber()
	if err != nil || latest == nil {
		return
	}
	if from == nil {
		from = latest
	}
	if to == nil {
		to = latest
	}
	return
}

func parseBlock(block string) (*utils.Uint256, error) {
	switch block {
	case "", "pending", "latest":
		return nil, nil
	case "earliest":
		zero := utils.IntToUint256(0)
		return &zero, nil
	default:
		val := utils.IntToUint256(0)
		err := val.FromHexString(block)
		if err != nil {
			return nil, err
		}
		return &val, nil
	}
}
