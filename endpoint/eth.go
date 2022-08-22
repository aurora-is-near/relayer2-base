package endpoint

import (
	"aurora-relayer-go-common/utils"
	"context"
)

var (
	protocolVersion = utils.IntToUint256(0x41)
	hashrate        = utils.IntToUint256(0)
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

func (e *Eth) Coinbase(_ context.Context) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_coinbase"); err != nil {
		return nil, err
	}
	return &hashrate, nil
}

func (e *Eth) ProtocolVersion(_ context.Context) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_protocolVersion"); err != nil {
		return nil, err
	}
	return &protocolVersion, nil
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
	return bn, nil
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
	count := utils.IntToUint256(*cnt)
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
	count := utils.IntToUint256(*cnt)
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
	idx := index.Int64()
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
	idx := index.Int64()
	tx, err := (*e.DbHandler).GetTransactionByBlockNumberAndIndex(number, idx)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return tx.ToResponse(), nil
}

// GetLogs returns an array of log objects for the given filter
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure or number not found, returns error code '-32000' with custom message.
func (e *Eth) GetLogs(_ context.Context, rawFilter utils.FilterOptions) (*[]utils.LogResponse, error) {
	if err := e.IsEndpointAllowed("eth_getLogs"); err != nil {
		return nil, err
	}
	filter, err := e.formatFilterOptions(&rawFilter)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	logResponses, err := (*e.DbHandler).GetLogs(filter)
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	return logResponses, nil
}

// NewFilter creates a new filter based on the filter options and returns newly created filter ID on success.
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// FilterOptions object is mandatory but all keys are optional
//
// "fromBlock": QUANTITY|TAG - Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
//
// "toBlock": QUANTITY|TAG - Integer block number, or "latest" for the last mined block or "pending", "earliest" for not yet mined transactions.
//
// "address": DATA|Array - Contract address or a list of addresses from which logs should originate.
//
// "topics": Array of DATA - Topics are order-dependent. Each topic can also be an array of DATA with "or" options.
//
// On filter option parsing failure, returns error code '32602' with custom message.
//
// On DB failure, returns error code '-32000' with custom message.
func (e *Eth) NewFilter(_ context.Context, filterOptions utils.FilterOptions) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_newFilter"); err != nil {
		return nil, err
	}

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

	_ = utils.StoredFilter{
		Type:      "event",
		CreatedBy: "0.0.0.0",
		BlockHash: parsedFilter.BlockHash,
		FromBlock: parsedFilter.FromBlock,
		ToBlock:   parsedFilter.ToBlock,
		Addresses: parsedFilter.Address,
		Topics:    parsedFilter.Topics,
		PollBlock: *bn,
	}

	// TODO store filter
	return filterId, nil
}

// NewBlockFilter creates a filter and returns newly created filter ID on success. To check if the state has changed,
// call "eth_getFilterChanges".
//
// If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
//
// On DB failure, returns error code '-32000' with custom message.
func (e *Eth) NewBlockFilter(_ context.Context) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_newBlockFilter"); err != nil {
		return nil, err
	}

	bn, err := (*e.DbHandler).BlockNumber()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}
	filterId, err := utils.RandomUint256()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}

	_ = &utils.StoredFilter{
		Type:      "block",
		CreatedBy: "0.0.0.0",
		PollBlock: *bn.Add(1),
	}

	// TODO store filter
	return filterId, nil
}

// UninstallFilter Uninstalls a filter with given filter id. Additionally, filters timeout when they aren't requested
// with "eth_getFilterChanges" for a period of time.
//
// TODO comment
func (e *Eth) UninstallFilter(_ context.Context, filterId utils.Uint256) (bool, error) {
	// TODO delete filter
	return false, nil
}

// GetFilterChanges polls method for a filter, on success returns an array of logs which occurred since last poll.
//
// TODO comment
func (e *Eth) GetFilterChanges(_ context.Context, filterId utils.Uint256) (*[]interface{}, error) {
	if err := e.IsEndpointAllowed("eth_getFilterChanges"); err != nil {
		return nil, err
	}

	bn, err := (*e.DbHandler).BlockNumber()
	if err != nil {
		return nil, &utils.GenericError{Err: err}
	}

	// TODO fetch filter
	var storedFilter utils.StoredFilter
	filterChanges := make([]interface{}, 0)

	// TODO see comment on parseFromAndTo
	switch storedFilter.Type {
	case "event":
		logFilter := utils.LogFilter{
			FromBlock: storedFilter.FromBlock,
			ToBlock:   storedFilter.ToBlock,
			Address:   storedFilter.Addresses,
			Topics:    storedFilter.Topics,
		}
		logs, err := (*e.DbHandler).GetLogs(&logFilter)
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
		// TODO specify error
	}

	storedFilter.PollBlock = *bn.Add(1)

	// TODO store filter

	return &filterChanges, nil
}

func (e *Eth) formatFilterOptions(filterOptions *utils.FilterOptions) (*utils.LogFilter, error) {
	result := &utils.LogFilter{}
	if filterOptions.BlockHash != nil {
		number, err := e.blockHashToNumber(*filterOptions.BlockHash)
		if err != nil {
			return nil, err
		}
		result.FromBlock = number
		result.ToBlock = number.Add(1)
	} else {
		var err error
		result.FromBlock, result.ToBlock, err = e.parseFromAndTo(filterOptions)
		if err != nil {
			return nil, err
		}
	}
	if filterOptions.Address != nil {
		addresses := parseAddresses(filterOptions.Address)
		result.Address = make(map[utils.Address]bool)
		for _, a := range addresses {
			result.Address[a] = true
		}
	}
	return result, nil
}

func (e *Eth) blockHashToNumber(hash utils.H256) (*utils.Uint256, error) {
	block, err := (*e.DbHandler).GetBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	return &block.Sequence, nil
}

// TODO: ask the correct way?
// => getLogChanges should return changes since the last poll, if a filter has from/to field it does it overwrite the
// 'since last change' condition?
//
// if it overwrites, for filter options of having empty from/to fields, the returned from/to fields must be empty as well
// otherwise we cannot differentiate that filter(1) (empty to/from field on purpose) from a non empty to/from option filter(2)
// if filter is of type (1) the from/to fields can be set during getLogs so that the "since last change" condition can be meet
// if filter is of type (2) then we already have from/to...
func (e *Eth) parseFromAndTo(filter *utils.FilterOptions) (from, to *utils.Uint256, err error) {
	fromBlock := parseBlock(filter.FromBlock)
	if fromBlock != nil {
		from = fromBlock
	}
	toBlock := parseBlock(filter.ToBlock)
	if toBlock != nil {
		to = toBlock.Add(1)
	}
	if from != nil && to != nil {
		return
	}

	latest, err := (*e.DbHandler).BlockNumber()
	if err != nil {
		return
	}
	if from == nil {
		from = latest
	}
	if to == nil {
		to = latest.Add(1)
	}
	return
}

func parseBlock(block interface{}) *utils.Uint256 {
	// TODO validation
	switch v := block.(type) {
	case string:
		switch v {
		case "", "pending", "latest":
			return nil
		case "earliest":
			zero := utils.IntToUint256(0)
			return &zero
		default:
			val := utils.IntToUint256(0)
			err := val.FromHexString(v)
			if err != nil {
				return nil
			}
			return &val
		}
	default:
		return nil
	}
}

func parseAddresses(addresses interface{}) []utils.Address {
	// TODO validation
	switch v := addresses.(type) {
	case []string:
		res := make([]utils.Address, 0, len(v))
		for _, val := range v {
			res = append(res, utils.HexStringToAddress(val))
		}
		return res

	case string:
		return []utils.Address{utils.HexStringToAddress(v)}
	default:
		return nil
	}
}
