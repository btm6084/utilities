package permutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	t.Run("2s", func(t *testing.T) {
		input := []string{"a", "b"}
		expected := [][]string{
			{"a", "b"},
			{"b", "a"},
		}

		actual := Strings(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 2*1, len(actual))
	})

	t.Run("3s", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		expected := [][]string{
			{"a", "b", "c"},
			{"b", "a", "c"},

			{"a", "c", "b"},
			{"b", "c", "a"},

			{"c", "a", "b"},
			{"c", "b", "a"},
		}

		actual := Strings(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 3*2*1, len(actual))
	})

	t.Run("4s", func(t *testing.T) {
		input := []string{"a", "b", "c", "d"}
		expected := [][]string{
			{"a", "b", "c", "d"},
			{"a", "c", "b", "d"},
			{"b", "a", "c", "d"},
			{"b", "c", "a", "d"},
			{"c", "b", "a", "d"},
			{"c", "a", "b", "d"},

			{"a", "b", "d", "c"},
			{"a", "c", "d", "b"},
			{"b", "a", "d", "c"},
			{"b", "c", "d", "a"},
			{"c", "b", "d", "a"},
			{"c", "a", "d", "b"},

			{"a", "d", "b", "c"},
			{"a", "d", "c", "b"},
			{"b", "d", "a", "c"},
			{"b", "d", "c", "a"},
			{"c", "d", "b", "a"},
			{"c", "d", "a", "b"},

			{"d", "a", "b", "c"},
			{"d", "a", "c", "b"},
			{"d", "b", "a", "c"},
			{"d", "b", "c", "a"},
			{"d", "c", "b", "a"},
			{"d", "c", "a", "b"},
		}

		actual := Strings(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 4*3*2*1, len(actual))
	})

	t.Run("5s", func(t *testing.T) {
		input := []string{"a", "b", "c", "d", "e"}
		expected := [][]string{
			{"a", "b", "c", "d", "e"},
			{"a", "c", "b", "d", "e"},
			{"b", "a", "c", "d", "e"},
			{"b", "c", "a", "d", "e"},
			{"c", "b", "a", "d", "e"},
			{"c", "a", "b", "d", "e"},
			{"a", "b", "d", "c", "e"},
			{"a", "c", "d", "b", "e"},
			{"b", "a", "d", "c", "e"},
			{"b", "c", "d", "a", "e"},
			{"c", "b", "d", "a", "e"},
			{"c", "a", "d", "b", "e"},
			{"a", "d", "b", "c", "e"},
			{"a", "d", "c", "b", "e"},
			{"b", "d", "a", "c", "e"},
			{"b", "d", "c", "a", "e"},
			{"c", "d", "b", "a", "e"},
			{"c", "d", "a", "b", "e"},
			{"d", "a", "b", "c", "e"},
			{"d", "a", "c", "b", "e"},
			{"d", "b", "a", "c", "e"},
			{"d", "b", "c", "a", "e"},
			{"d", "c", "b", "a", "e"},
			{"d", "c", "a", "b", "e"},

			{"a", "b", "c", "e", "d"},
			{"a", "c", "b", "e", "d"},
			{"b", "a", "c", "e", "d"},
			{"b", "c", "a", "e", "d"},
			{"c", "b", "a", "e", "d"},
			{"c", "a", "b", "e", "d"},
			{"a", "b", "d", "e", "c"},
			{"a", "c", "d", "e", "b"},
			{"b", "a", "d", "e", "c"},
			{"b", "c", "d", "e", "a"},
			{"c", "b", "d", "e", "a"},
			{"c", "a", "d", "e", "b"},
			{"a", "d", "b", "e", "c"},
			{"a", "d", "c", "e", "b"},
			{"b", "d", "a", "e", "c"},
			{"b", "d", "c", "e", "a"},
			{"c", "d", "b", "e", "a"},
			{"c", "d", "a", "e", "b"},
			{"d", "a", "b", "e", "c"},
			{"d", "a", "c", "e", "b"},
			{"d", "b", "a", "e", "c"},
			{"d", "b", "c", "e", "a"},
			{"d", "c", "b", "e", "a"},
			{"d", "c", "a", "e", "b"},

			{"a", "b", "e", "c", "d"},
			{"a", "c", "e", "b", "d"},
			{"b", "a", "e", "c", "d"},
			{"b", "c", "e", "a", "d"},
			{"c", "b", "e", "a", "d"},
			{"c", "a", "e", "b", "d"},
			{"a", "b", "e", "d", "c"},
			{"a", "c", "e", "d", "b"},
			{"b", "a", "e", "d", "c"},
			{"b", "c", "e", "d", "a"},
			{"c", "b", "e", "d", "a"},
			{"c", "a", "e", "d", "b"},
			{"a", "d", "e", "b", "c"},
			{"a", "d", "e", "c", "b"},
			{"b", "d", "e", "a", "c"},
			{"b", "d", "e", "c", "a"},
			{"c", "d", "e", "b", "a"},
			{"c", "d", "e", "a", "b"},
			{"d", "a", "e", "b", "c"},
			{"d", "a", "e", "c", "b"},
			{"d", "b", "e", "a", "c"},
			{"d", "b", "e", "c", "a"},
			{"d", "c", "e", "b", "a"},
			{"d", "c", "e", "a", "b"},

			{"a", "e", "b", "c", "d"},
			{"a", "e", "c", "b", "d"},
			{"b", "e", "a", "c", "d"},
			{"b", "e", "c", "a", "d"},
			{"c", "e", "b", "a", "d"},
			{"c", "e", "a", "b", "d"},
			{"a", "e", "b", "d", "c"},
			{"a", "e", "c", "d", "b"},
			{"b", "e", "a", "d", "c"},
			{"b", "e", "c", "d", "a"},
			{"c", "e", "b", "d", "a"},
			{"c", "e", "a", "d", "b"},
			{"a", "e", "d", "b", "c"},
			{"a", "e", "d", "c", "b"},
			{"b", "e", "d", "a", "c"},
			{"b", "e", "d", "c", "a"},
			{"c", "e", "d", "b", "a"},
			{"c", "e", "d", "a", "b"},
			{"d", "e", "a", "b", "c"},
			{"d", "e", "a", "c", "b"},
			{"d", "e", "b", "a", "c"},
			{"d", "e", "b", "c", "a"},
			{"d", "e", "c", "b", "a"},
			{"d", "e", "c", "a", "b"},

			{"e", "a", "b", "c", "d"},
			{"e", "a", "c", "b", "d"},
			{"e", "b", "a", "c", "d"},
			{"e", "b", "c", "a", "d"},
			{"e", "c", "b", "a", "d"},
			{"e", "c", "a", "b", "d"},
			{"e", "a", "b", "d", "c"},
			{"e", "a", "c", "d", "b"},
			{"e", "b", "a", "d", "c"},
			{"e", "b", "c", "d", "a"},
			{"e", "c", "b", "d", "a"},
			{"e", "c", "a", "d", "b"},
			{"e", "a", "d", "b", "c"},
			{"e", "a", "d", "c", "b"},
			{"e", "b", "d", "a", "c"},
			{"e", "b", "d", "c", "a"},
			{"e", "c", "d", "b", "a"},
			{"e", "c", "d", "a", "b"},
			{"e", "d", "a", "b", "c"},
			{"e", "d", "a", "c", "b"},
			{"e", "d", "b", "a", "c"},
			{"e", "d", "b", "c", "a"},
			{"e", "d", "c", "b", "a"},
			{"e", "d", "c", "a", "b"},
		}

		actual := Strings(input)
		assert.ElementsMatch(t, expected, actual)
		assert.Equal(t, 5*4*3*2*1, len(actual))
	})
}
