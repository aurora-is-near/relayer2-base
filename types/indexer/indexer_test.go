package indexer

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"

	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

func TestTopicMarshalJSON(t *testing.T) {
	topic := Topic(primitives.Data32FromHex("0x1234"))
	res, err := jsoniter.Marshal(topic)
	require.NoError(t, err)
	require.EqualValues(t, `"0x0000000000000000000000000000000000000000000000000000000000001234"`, string(res))
}
