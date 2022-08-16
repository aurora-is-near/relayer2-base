package endpoint

import (
	"aurora-relayer-go-common/utils"
	"context"
	"fmt"
)

type Eth struct {
	*Endpoint
}

func NewEth(endpoint *Endpoint) *Eth {
	return &Eth{endpoint}
}

func (e *Eth) ProtocolVersion(_ context.Context) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_protocolVersion"); err != nil {
		return nil, err
	}
	v := utils.Uint256("41")
	return &v, nil
}

func (e *Eth) Hashrate(_ context.Context) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_hashrate"); err != nil {
		return nil, err
	}
	v := utils.Uint256("0")
	return &v, nil
}

func (e *Eth) BlockNumber(_ context.Context) (*uint64, error) {
	if err := e.IsEndpointAllowed("eth_blockNumber"); err != nil {
		return nil, err
	}
	return (*e.DbHandler).BlockNumber()
}

func (e *Eth) GetBlockByHash(_ context.Context, hash utils.H256, isFull bool) (*utils.BlockResponse, error) {
	if err := e.IsEndpointAllowed("eth_getBlockByHash"); err != nil {
		return nil, err
	}
	// TODO implement
	return nil, nil
}

func (e *Eth) GetBlockByNumber(_ context.Context, number utils.Uint256, isFull bool) (*utils.BlockResponse, error) {
	if err := e.IsEndpointAllowed("eth_getBlockByNumber"); err != nil {
		return nil, err
	}
	// TODO implement
	return nil, nil
}

func (e *Eth) GetBlockTransactionCountByHash(_ context.Context, hash utils.H256) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_getBlockTransactionCountByHash"); err != nil {
		return nil, err
	}
	// TODO implement
	return nil, nil
}

func (e *Eth) GetBlockTransactionCountByNumber(_ context.Context, number utils.Uint256) (*utils.Uint256, error) {
	if err := e.IsEndpointAllowed("eth_getBlockTransactionCountByNumber"); err != nil {
		return nil, err
	}
	// TODO implement
	return nil, nil
}

func (e *Eth) GetTransactionByHash(_ context.Context, hash utils.H256) (*utils.TransactionResponse, error) {
	if err := e.IsEndpointAllowed("eth_getTransactionByHash"); err != nil {
		return nil, err
	}
	tx, err := (*e.DbHandler).GetTransactionByHash(hash)
	if err != nil {
		return nil, err
	}
	return tx.ToResponse(), nil
}

func (e *Eth) GetTransactionByBlockHashAndIndex(_ context.Context, hash utils.H256, index utils.Uint256) (*utils.TransactionResponse, error) {
	if err := e.IsEndpointAllowed("eth_getTransactionByBlockHashAndIndex"); err != nil {
		return nil, err
	}
	// TODO implement
	return nil, nil
}

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
		return nil, err
	}
	return tx.ToResponse(), nil
}

func (e *Eth) GetLogs(addr utils.Address, bn utils.Uint256, topic ...[]string) (*utils.Log, error) {

	// TODO implement
	return nil, nil
}
