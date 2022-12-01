package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSum(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		in := []int{1, 2, 3, 4, 5, 99999}

		require.Equal(t, 100014, Sum(in...))
	})

	t.Run("string", func(t *testing.T) {
		in := []string{"1", "2", "99999", "3", "4", "5"}
		require.Equal(t, "1299999345", Sum(in...))
	})

	t.Run("float", func(t *testing.T) {
		in := []float64{1.1, 2.2, 5.1, 3, 4.4, 5, 5.01}

		require.Equal(t, 25.810000000000002, Sum(in...))
	})
}

func TestMax(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		in := []int{1, 2, 3, 4, 5, 99999}

		require.Equal(t, 99999, Max(in...))
	})

	t.Run("string", func(t *testing.T) {
		in := []string{"1", "2", "99999", "3", "4", "5"}
		require.Equal(t, "99999", Max(in...))
	})

	t.Run("float", func(t *testing.T) {
		in := []float64{1.1, 2.2, 5.1, 3, 4.4, 5, 5.01}

		require.Equal(t, 5.1, Max(in...))
	})
}

func TestContains(t *testing.T) {

	t.Run("int", func(t *testing.T) {
		in := []int{1, 2, 3, 4, 5, 99999}

		require.Equal(t, true, Contains(in, 1))
		require.Equal(t, true, Contains(in, 99999))
		require.Equal(t, false, Contains(in, 123))
	})

	t.Run("string", func(t *testing.T) {
		in := []string{"1", "2", "3", "4", "5", "99999"}

		require.Equal(t, true, Contains(in, "1"))
		require.Equal(t, true, Contains(in, "99999"))
		require.Equal(t, false, Contains(in, "123"))
	})

}

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
