package stack

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrace(t *testing.T) {
	f, l := Trace(0)
	require.True(t, strings.HasSuffix(f, "utilities/stack/stack_test.go"), f)
	require.Equal(t, 11, l)

	func() {
		f, l = Trace(0)
		require.True(t, strings.HasSuffix(f, "utilities/stack/stack_test.go"), f)
		require.Equal(t, 16, l)
	}()

	f, l = trace()
	require.True(t, strings.HasSuffix(f, "utilities/stack/stack_test.go"), f)
	require.Equal(t, 21, l)
}

func trace() (string, int) {
	return Trace(1)
}
