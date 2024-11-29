package request

import (
	"errors"
	"fmt"

	"github.com/aurora-is-near/relayer2-base/types/common"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	jsoniter "github.com/json-iterator/go"
)

type Topics [][]primitives.Data32
type SingleOrSliceOfAddress []common.Address

type LogSubscriptionOptions struct {
	Address []common.Address `json:"address"`
	Topics  Topics           `json:"topics"`
}

type Filter struct {
	BlockHash *common.H256           `json:"blockhash"`
	FromBlock *common.BN64           `json:"fromBlock"`
	ToBlock   *common.BN64           `json:"toBlock"`
	Addresses SingleOrSliceOfAddress `json:"address"`
	Topics    Topics                 `json:"topics"`
}

func (t *Topics) UnmarshalJSON(b []byte) error {
	tps := [4]interface{}{}
	err := jsoniter.Unmarshal(b, &tps)
	if err != nil {
		return err
	}
	results := Topics{{}, {}, {}, {}}
	for i, t := range tps {
		switch v := t.(type) {
		case string:
			data, err := primitives.Data32FromHex(v)
			if err != nil {
				return err
			}

			results[i] = append(results[i], data)
		case []interface{}:
			for _, topic := range v {
				if topic, ok := topic.(string); ok {
					data, err := primitives.Data32FromHex(topic)
					if err != nil {
						return err
					}

					results[i] = append(results[i], data)
				}
			}
		case nil:
		default:
		}
	}
	*t = results
	return nil
}

func (a *SingleOrSliceOfAddress) UnmarshalJSON(b []byte) error {
	type input interface{}
	var raw input
	if err := jsoniter.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw != nil {
		// rawAddr can contain a single address or an slice of addresses
		switch rawAddr := raw.(type) {
		case string:
			addr, err := common.HexStringToAddress(rawAddr)
			if err != nil {
				return err
			}

			*a = []common.Address{addr}
		case []interface{}:
			for i, addr := range rawAddr {
				if strAddr, ok := addr.(string); ok {
					addr, err := common.HexStringToAddress(strAddr)
					if err != nil {
						return err
					}

					*a = append(*a, addr)
				} else {
					return fmt.Errorf("non-string address at index %d", i)
				}
			}
		default:
			return errors.New("invalid addresses field in filter options object")
		}
	}

	return nil
}
