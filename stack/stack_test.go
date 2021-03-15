package stack

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	f, l := Trace(0)
	assert.True(t, strings.HasSuffix(f, "utilities/stack/stack_test.go"), f)
	assert.Equal(t, 11, l)

	func() {
		f, l = Trace(0)
		assert.True(t, strings.HasSuffix(f, "utilities/stack/stack_test.go"), f)
		assert.Equal(t, 16, l)
	}()

	f, l = trace()
	assert.True(t, strings.HasSuffix(f, "utilities/stack/stack_test.go"), f)
	assert.Equal(t, 21, l)
}

func trace() (string, int) {
	return Trace(1)
}
