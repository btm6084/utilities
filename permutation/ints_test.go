package permutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInts(t *testing.T) {
	t.Run("2s", func(t *testing.T) {
		input := []int{1, 2}
		expected := [][]int{
			{1, 2},
			{2, 1},
		}

		actual := Ints(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 2*1, len(actual))
	})

	t.Run("3s", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := [][]int{
			{1, 2, 3},
			{2, 1, 3},

			{1, 3, 2},
			{2, 3, 1},

			{3, 1, 2},
			{3, 2, 1},
		}

		actual := Ints(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 3*2*1, len(actual))
	})

	t.Run("4s", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		expected := [][]int{
			{1, 2, 3, 4},
			{1, 3, 2, 4},
			{2, 1, 3, 4},
			{2, 3, 1, 4},
			{3, 2, 1, 4},
			{3, 1, 2, 4},

			{1, 2, 4, 3},
			{1, 3, 4, 2},
			{2, 1, 4, 3},
			{2, 3, 4, 1},
			{3, 2, 4, 1},
			{3, 1, 4, 2},

			{1, 4, 2, 3},
			{1, 4, 3, 2},
			{2, 4, 1, 3},
			{2, 4, 3, 1},
			{3, 4, 2, 1},
			{3, 4, 1, 2},

			{4, 1, 2, 3},
			{4, 1, 3, 2},
			{4, 2, 1, 3},
			{4, 2, 3, 1},
			{4, 3, 2, 1},
			{4, 3, 1, 2},
		}

		actual := Ints(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 4*3*2*1, len(actual))
	})

	t.Run("5s", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := [][]int{
			{1, 2, 3, 4, 5},
			{1, 3, 2, 4, 5},
			{2, 1, 3, 4, 5},
			{2, 3, 1, 4, 5},
			{3, 2, 1, 4, 5},
			{3, 1, 2, 4, 5},
			{1, 2, 4, 3, 5},
			{1, 3, 4, 2, 5},
			{2, 1, 4, 3, 5},
			{2, 3, 4, 1, 5},
			{3, 2, 4, 1, 5},
			{3, 1, 4, 2, 5},
			{1, 4, 2, 3, 5},
			{1, 4, 3, 2, 5},
			{2, 4, 1, 3, 5},
			{2, 4, 3, 1, 5},
			{3, 4, 2, 1, 5},
			{3, 4, 1, 2, 5},
			{4, 1, 2, 3, 5},
			{4, 1, 3, 2, 5},
			{4, 2, 1, 3, 5},
			{4, 2, 3, 1, 5},
			{4, 3, 2, 1, 5},
			{4, 3, 1, 2, 5},

			{1, 2, 3, 5, 4},
			{1, 3, 2, 5, 4},
			{2, 1, 3, 5, 4},
			{2, 3, 1, 5, 4},
			{3, 2, 1, 5, 4},
			{3, 1, 2, 5, 4},
			{1, 2, 4, 5, 3},
			{1, 3, 4, 5, 2},
			{2, 1, 4, 5, 3},
			{2, 3, 4, 5, 1},
			{3, 2, 4, 5, 1},
			{3, 1, 4, 5, 2},
			{1, 4, 2, 5, 3},
			{1, 4, 3, 5, 2},
			{2, 4, 1, 5, 3},
			{2, 4, 3, 5, 1},
			{3, 4, 2, 5, 1},
			{3, 4, 1, 5, 2},
			{4, 1, 2, 5, 3},
			{4, 1, 3, 5, 2},
			{4, 2, 1, 5, 3},
			{4, 2, 3, 5, 1},
			{4, 3, 2, 5, 1},
			{4, 3, 1, 5, 2},

			{1, 2, 5, 3, 4},
			{1, 3, 5, 2, 4},
			{2, 1, 5, 3, 4},
			{2, 3, 5, 1, 4},
			{3, 2, 5, 1, 4},
			{3, 1, 5, 2, 4},
			{1, 2, 5, 4, 3},
			{1, 3, 5, 4, 2},
			{2, 1, 5, 4, 3},
			{2, 3, 5, 4, 1},
			{3, 2, 5, 4, 1},
			{3, 1, 5, 4, 2},
			{1, 4, 5, 2, 3},
			{1, 4, 5, 3, 2},
			{2, 4, 5, 1, 3},
			{2, 4, 5, 3, 1},
			{3, 4, 5, 2, 1},
			{3, 4, 5, 1, 2},
			{4, 1, 5, 2, 3},
			{4, 1, 5, 3, 2},
			{4, 2, 5, 1, 3},
			{4, 2, 5, 3, 1},
			{4, 3, 5, 2, 1},
			{4, 3, 5, 1, 2},

			{1, 5, 2, 3, 4},
			{1, 5, 3, 2, 4},
			{2, 5, 1, 3, 4},
			{2, 5, 3, 1, 4},
			{3, 5, 2, 1, 4},
			{3, 5, 1, 2, 4},
			{1, 5, 2, 4, 3},
			{1, 5, 3, 4, 2},
			{2, 5, 1, 4, 3},
			{2, 5, 3, 4, 1},
			{3, 5, 2, 4, 1},
			{3, 5, 1, 4, 2},
			{1, 5, 4, 2, 3},
			{1, 5, 4, 3, 2},
			{2, 5, 4, 1, 3},
			{2, 5, 4, 3, 1},
			{3, 5, 4, 2, 1},
			{3, 5, 4, 1, 2},
			{4, 5, 1, 2, 3},
			{4, 5, 1, 3, 2},
			{4, 5, 2, 1, 3},
			{4, 5, 2, 3, 1},
			{4, 5, 3, 2, 1},
			{4, 5, 3, 1, 2},

			{5, 1, 2, 3, 4},
			{5, 1, 3, 2, 4},
			{5, 2, 1, 3, 4},
			{5, 2, 3, 1, 4},
			{5, 3, 2, 1, 4},
			{5, 3, 1, 2, 4},
			{5, 1, 2, 4, 3},
			{5, 1, 3, 4, 2},
			{5, 2, 1, 4, 3},
			{5, 2, 3, 4, 1},
			{5, 3, 2, 4, 1},
			{5, 3, 1, 4, 2},
			{5, 1, 4, 2, 3},
			{5, 1, 4, 3, 2},
			{5, 2, 4, 1, 3},
			{5, 2, 4, 3, 1},
			{5, 3, 4, 2, 1},
			{5, 3, 4, 1, 2},
			{5, 4, 1, 2, 3},
			{5, 4, 1, 3, 2},
			{5, 4, 2, 1, 3},
			{5, 4, 2, 3, 1},
			{5, 4, 3, 2, 1},
			{5, 4, 3, 1, 2},
		}

		actual := Ints(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 5*4*3*2*1, len(actual))
	})
}

func TestIntsRecursive(t *testing.T) {
	t.Run("2s", func(t *testing.T) {
		input := []int{1, 2}
		expected := [][]int{
			{1, 2},
			{2, 1},
		}

		actual := IntsRecursive(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 2*1, len(actual))
	})

	t.Run("3s", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := [][]int{
			{1, 2, 3},
			{2, 1, 3},

			{1, 3, 2},
			{2, 3, 1},

			{3, 1, 2},
			{3, 2, 1},
		}

		actual := IntsRecursive(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 3*2*1, len(actual))
	})

	t.Run("4s", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		expected := [][]int{
			{1, 2, 3, 4},
			{1, 3, 2, 4},
			{2, 1, 3, 4},
			{2, 3, 1, 4},
			{3, 2, 1, 4},
			{3, 1, 2, 4},

			{1, 2, 4, 3},
			{1, 3, 4, 2},
			{2, 1, 4, 3},
			{2, 3, 4, 1},
			{3, 2, 4, 1},
			{3, 1, 4, 2},

			{1, 4, 2, 3},
			{1, 4, 3, 2},
			{2, 4, 1, 3},
			{2, 4, 3, 1},
			{3, 4, 2, 1},
			{3, 4, 1, 2},

			{4, 1, 2, 3},
			{4, 1, 3, 2},
			{4, 2, 1, 3},
			{4, 2, 3, 1},
			{4, 3, 2, 1},
			{4, 3, 1, 2},
		}

		actual := IntsRecursive(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 4*3*2*1, len(actual))
	})

	t.Run("5s", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := [][]int{
			{1, 2, 3, 4, 5},
			{1, 3, 2, 4, 5},
			{2, 1, 3, 4, 5},
			{2, 3, 1, 4, 5},
			{3, 2, 1, 4, 5},
			{3, 1, 2, 4, 5},
			{1, 2, 4, 3, 5},
			{1, 3, 4, 2, 5},
			{2, 1, 4, 3, 5},
			{2, 3, 4, 1, 5},
			{3, 2, 4, 1, 5},
			{3, 1, 4, 2, 5},
			{1, 4, 2, 3, 5},
			{1, 4, 3, 2, 5},
			{2, 4, 1, 3, 5},
			{2, 4, 3, 1, 5},
			{3, 4, 2, 1, 5},
			{3, 4, 1, 2, 5},
			{4, 1, 2, 3, 5},
			{4, 1, 3, 2, 5},
			{4, 2, 1, 3, 5},
			{4, 2, 3, 1, 5},
			{4, 3, 2, 1, 5},
			{4, 3, 1, 2, 5},

			{1, 2, 3, 5, 4},
			{1, 3, 2, 5, 4},
			{2, 1, 3, 5, 4},
			{2, 3, 1, 5, 4},
			{3, 2, 1, 5, 4},
			{3, 1, 2, 5, 4},
			{1, 2, 4, 5, 3},
			{1, 3, 4, 5, 2},
			{2, 1, 4, 5, 3},
			{2, 3, 4, 5, 1},
			{3, 2, 4, 5, 1},
			{3, 1, 4, 5, 2},
			{1, 4, 2, 5, 3},
			{1, 4, 3, 5, 2},
			{2, 4, 1, 5, 3},
			{2, 4, 3, 5, 1},
			{3, 4, 2, 5, 1},
			{3, 4, 1, 5, 2},
			{4, 1, 2, 5, 3},
			{4, 1, 3, 5, 2},
			{4, 2, 1, 5, 3},
			{4, 2, 3, 5, 1},
			{4, 3, 2, 5, 1},
			{4, 3, 1, 5, 2},

			{1, 2, 5, 3, 4},
			{1, 3, 5, 2, 4},
			{2, 1, 5, 3, 4},
			{2, 3, 5, 1, 4},
			{3, 2, 5, 1, 4},
			{3, 1, 5, 2, 4},
			{1, 2, 5, 4, 3},
			{1, 3, 5, 4, 2},
			{2, 1, 5, 4, 3},
			{2, 3, 5, 4, 1},
			{3, 2, 5, 4, 1},
			{3, 1, 5, 4, 2},
			{1, 4, 5, 2, 3},
			{1, 4, 5, 3, 2},
			{2, 4, 5, 1, 3},
			{2, 4, 5, 3, 1},
			{3, 4, 5, 2, 1},
			{3, 4, 5, 1, 2},
			{4, 1, 5, 2, 3},
			{4, 1, 5, 3, 2},
			{4, 2, 5, 1, 3},
			{4, 2, 5, 3, 1},
			{4, 3, 5, 2, 1},
			{4, 3, 5, 1, 2},

			{1, 5, 2, 3, 4},
			{1, 5, 3, 2, 4},
			{2, 5, 1, 3, 4},
			{2, 5, 3, 1, 4},
			{3, 5, 2, 1, 4},
			{3, 5, 1, 2, 4},
			{1, 5, 2, 4, 3},
			{1, 5, 3, 4, 2},
			{2, 5, 1, 4, 3},
			{2, 5, 3, 4, 1},
			{3, 5, 2, 4, 1},
			{3, 5, 1, 4, 2},
			{1, 5, 4, 2, 3},
			{1, 5, 4, 3, 2},
			{2, 5, 4, 1, 3},
			{2, 5, 4, 3, 1},
			{3, 5, 4, 2, 1},
			{3, 5, 4, 1, 2},
			{4, 5, 1, 2, 3},
			{4, 5, 1, 3, 2},
			{4, 5, 2, 1, 3},
			{4, 5, 2, 3, 1},
			{4, 5, 3, 2, 1},
			{4, 5, 3, 1, 2},

			{5, 1, 2, 3, 4},
			{5, 1, 3, 2, 4},
			{5, 2, 1, 3, 4},
			{5, 2, 3, 1, 4},
			{5, 3, 2, 1, 4},
			{5, 3, 1, 2, 4},
			{5, 1, 2, 4, 3},
			{5, 1, 3, 4, 2},
			{5, 2, 1, 4, 3},
			{5, 2, 3, 4, 1},
			{5, 3, 2, 4, 1},
			{5, 3, 1, 4, 2},
			{5, 1, 4, 2, 3},
			{5, 1, 4, 3, 2},
			{5, 2, 4, 1, 3},
			{5, 2, 4, 3, 1},
			{5, 3, 4, 2, 1},
			{5, 3, 4, 1, 2},
			{5, 4, 1, 2, 3},
			{5, 4, 1, 3, 2},
			{5, 4, 2, 1, 3},
			{5, 4, 2, 3, 1},
			{5, 4, 3, 2, 1},
			{5, 4, 3, 1, 2},
		}

		actual := IntsRecursive(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 5*4*3*2*1, len(actual))
	})
}
