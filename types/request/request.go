package request

import (
	"aurora-relayer-go-common/types/common"
	"encoding/json"
)

type Topics [][][]byte

type LogSubscriptionOptions struct {
	Address []common.Address `json:"address"`
	Topics  Topics           `json:"topics"`
}

type Filter struct {
	BlockHash *common.H256     `json:"blockhash"`
	FromBlock *common.BN64     `json:"fromBlock"`
	ToBlock   *common.BN64     `json:"toBlock"`
	Addresses []common.Address `json:"address"`
	Topics    Topics           `json:"topics"`
}

func (t *Topics) UnmarshalJSON(b []byte) error {
	tps := [4]interface{}{}
	err := json.Unmarshal(b, &tps)
	if err != nil {
		return err
	}
	results := Topics{{}, {}, {}, {}}
	for i, t := range tps {
		switch v := t.(type) {
		case string:
			results[i] = append(results[i], []byte(v))
		case []interface{}:
			for _, topic := range v {
				if topic, ok := topic.(string); ok {
					results[i] = append(results[i], []byte(topic))
				}
			}
		case nil:
		default:
		}
	}
	*t = results
	return nil
}
