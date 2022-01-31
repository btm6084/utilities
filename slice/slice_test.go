package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntToString(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 99999}
	ex := []string{"1", "2", "3", "4", "5", "99999"}

	require.Equal(t, ex, IntToString(in))
}
