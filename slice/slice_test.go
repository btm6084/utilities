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

func TestUnique(t *testing.T) {
	t.Run("Int", func(t *testing.T) {
		in := []int{1, 1, 2, 18273, 18272, 3, 4, 2, 5, 99999, 18273}
		ex := []int{1, 2, 18273, 18272, 3, 4, 5, 99999}

		require.Equal(t, ex, Unique(in))
	})
	t.Run("String", func(t *testing.T) {
		in := []string{"1", "1", "2", "18273", "18272", "3", "4", "2", "5", "99999", "18273"}
		ex := []string{"1", "2", "18273", "18272", "3", "4", "5", "99999"}

		require.Equal(t, ex, Unique(in))
	})
}

func TestToInterface(t *testing.T) {
	t.Run("Int", func(t *testing.T) {
		in := []int{1, 1, 2, 18273, 18272, 3, 4, 2, 5, 99999, 18273}
		ex := []interface{}{1, 1, 2, 18273, 18272, 3, 4, 2, 5, 99999, 18273}

		require.Equal(t, ex, ToInterface(in))
	})
	t.Run("String", func(t *testing.T) {
		in := []string{"1", "1", "2", "18273", "18272", "3", "4", "2", "5", "99999", "18273"}
		ex := []interface{}{"1", "1", "2", "18273", "18272", "3", "4", "2", "5", "99999", "18273"}

		require.Equal(t, ex, ToInterface(in))
	})
}
