package dbschema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	path := Path(Const(1), Var(3), Var(1), Const(0), Const(0))

	result := path.Get(uint64(5*256*256+2*256+8), uint64(255))
	require.Equal(t, []byte{1, 5, 2, 8, 255, 0, 0}, result)
	require.True(t, path.Matches(result))
	require.Equal(t, []byte{5, 2, 8}, path.ReadVar(result, 0))
	require.Equal(t, uint64(5*256*256+2*256+8), path.ReadUintVar(result, 0))
	require.Equal(t, []byte{255}, path.ReadVar(result, 1))
	require.Equal(t, uint64(255), path.ReadUintVar(result, 1))

	result = path.Get(uint64(4*256*256+9*256+1), []byte{128})
	require.Equal(t, []byte{1, 4, 9, 1, 128, 0, 0}, result)
	require.True(t, path.Matches(result))
	require.Equal(t, []byte{4, 9, 1}, path.ReadVar(result, 0))
	require.Equal(t, uint64(4*256*256+9*256+1), path.ReadUintVar(result, 0))
	require.Equal(t, []byte{128}, path.ReadVar(result, 1))
	require.Equal(t, uint64(128), path.ReadUintVar(result, 1))

	result = path.Get([]byte{3, 14, 15}, uint64(0))
	require.Equal(t, []byte{1, 3, 14, 15, 0, 0, 0}, result)
	require.True(t, path.Matches(result))
	require.Equal(t, []byte{3, 14, 15}, path.ReadVar(result, 0))
	require.Equal(t, uint64(3*256*256+14*256+15), path.ReadUintVar(result, 0))
	require.Equal(t, []byte{0}, path.ReadVar(result, 1))
	require.Equal(t, uint64(0), path.ReadUintVar(result, 1))

	result = path.Get([]byte{0, 0, 0}, []byte{1})
	require.Equal(t, []byte{1, 0, 0, 0, 1, 0, 0}, result)
	require.True(t, path.Matches(result))
	require.Equal(t, []byte{0, 0, 0}, path.ReadVar(result, 0))
	require.Equal(t, uint64(0*256*256+0*256+0), path.ReadUintVar(result, 0))
	require.Equal(t, []byte{1}, path.ReadVar(result, 1))
	require.Equal(t, uint64(1), path.ReadUintVar(result, 1))

	require.False(t, path.Matches([]byte{}))
	require.False(t, path.Matches([]byte{0, 0, 0, 0, 1, 0, 0}))
	require.False(t, path.Matches([]byte{1, 0, 0, 0, 1, 1, 0}))
	require.False(t, path.Matches([]byte{1, 0, 0, 0, 1, 0, 255}))

	path = Path(Var(8))
	result = path.Get(uint64(12345678911223344))
	require.Equal(t, uint64(12345678911223344), path.ReadUintVar(result, 0))
	require.True(t, path.Matches([]byte{0, 1, 2, 3, 4, 5, 6, 7}))
	require.True(t, path.Matches([]byte{0, 0, 0, 0, 0, 0, 0, 0}))
	require.True(t, path.Matches([]byte{255, 255, 255, 255, 255, 255, 255, 255}))
	require.False(t, path.Matches([]byte{0, 1, 2, 3, 4, 5, 6}))
}
