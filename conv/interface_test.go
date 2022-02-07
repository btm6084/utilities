package conv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Thing interface {
	Stuff(s string)
}

func TestIsNil(t *testing.T) {
	require.True(t, IsNil((*string)(nil)))
	require.True(t, IsNil((*int)(nil)))
	require.True(t, IsNil((*float64)(nil)))
	require.True(t, IsNil((*struct{})(nil)))
	require.True(t, IsNil((interface{})(nil)))
	require.True(t, IsNil((*interface{})(nil)))
	require.True(t, IsNil((*Thing)(nil)))

	var thing Thing
	require.True(t, IsNil(thing))

	var empty []byte
	require.True(t, IsNil(empty))

	require.False(t, IsNil("string"))
	require.False(t, IsNil(22))
	require.False(t, IsNil(4.5))
	require.False(t, IsNil(struct{ A string }{A: "Hello"}))
	require.False(t, IsNil([]interface{}{"Hi"}))
	require.False(t, IsNil([]byte{}))
}
