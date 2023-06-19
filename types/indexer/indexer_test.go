package indexer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

func TestTopicMarshalJSON(t *testing.T) {
	topic := Topic(primitives.Data32FromHex("0x1234"))
	res, err := json.Marshal(topic)
	require.NoError(t, err)
	require.EqualValues(t, `"0x1234000000000000000000000000000000000000000000000000000000000000"`, string(res))
}
