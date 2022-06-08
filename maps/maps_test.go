package maps

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeys(t *testing.T) {
	t.Run("Strings", func(t *testing.T) {
		in := map[string]int{"a": 1, "b": 2, "c": 3}
		out := []string{"a", "b", "c"}

		actual := Keys(in)
		for _, v := range actual {
			require.Contains(t, out, v)
		}
	})

	t.Run("Ints", func(t *testing.T) {
		in := map[int]int{1: 1, 2: 2, 3: 3}
		out := []int{1, 2, 3}

		actual := Keys(in)
		for _, v := range actual {
			require.Contains(t, out, v)
		}
	})
}

func TestValues(t *testing.T) {
	t.Run("Ints", func(t *testing.T) {
		in := map[string]int{"a": 1, "b": 2, "c": 3}
		out := []int{1, 2, 3}

		actual := Values(in)
		for _, v := range actual {
			require.Contains(t, out, v)
		}
	})

	t.Run("Strings", func(t *testing.T) {
		in := map[int]string{1: "a", 2: "b", 3: "c"}
		out := []string{"a", "b", "c"}

		actual := Values(in)
		for _, v := range actual {
			require.Contains(t, out, v)
		}
	})
}
