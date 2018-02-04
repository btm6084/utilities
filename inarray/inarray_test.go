package inarray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	input := []string{"a", "b", "c", "Random Stuff"}

	assert.Equal(t, 0, Strings("a", input))
	assert.Equal(t, 1, Strings("b", input))
	assert.Equal(t, 2, Strings("c", input))
	assert.Equal(t, -1, Strings("d", input))
	assert.Equal(t, -1, Strings(" ", input))
	assert.Equal(t, -1, Strings("Random stuff", input))
	assert.Equal(t, 3, Strings("Random Stuff", input))
}

func TestInts(t *testing.T) {
	input := []int{19, 17, 1543}

	assert.Equal(t, 0, Ints(19, input))
	assert.Equal(t, 1, Ints(17, input))
	assert.Equal(t, 2, Ints(1543, input))
	assert.Equal(t, -1, Ints(2435, input))
	assert.Equal(t, -1, Ints(0, input))
	assert.Equal(t, -1, Ints(5, input))
}

func TestFloats(t *testing.T) {
	input := []float64{19.16, 17.99, 0.0001543}

	assert.Equal(t, 0, Floats(19.16, input))
	assert.Equal(t, 1, Floats(17.99, input))
	assert.Equal(t, 2, Floats(0.0001543, input))
	assert.Equal(t, -1, Floats(2435.0, input))
	assert.Equal(t, -1, Floats(0.0, input))
	assert.Equal(t, -1, Floats(5.0, input))
}
