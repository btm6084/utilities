package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenericStack(t *testing.T) {
	t.Run("Bool", func(t *testing.T) {
		s := Stack[bool]{}

		s.Push(true)
		s.Push(false)
		s.Push(true)
		s.Push(false)
		require.Equal(t, 4, s.Len())

		i, err := s.Pop()
		require.Equal(t, false, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, true, i)
		require.Nil(t, err)

		i, err = s.Pop()
		require.Equal(t, true, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, false, i)
		require.Nil(t, err)

		i, err = s.Pop()
		require.Equal(t, false, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, true, i)
		require.Nil(t, err)

		i, err = s.Pop()
		require.Equal(t, true, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, false, i)
		require.Equal(t, ErrStackEmpty, err)

		i, err = s.Pop()
		require.Equal(t, false, i)
		require.Equal(t, ErrStackEmpty, err)
	})
	t.Run("Float", func(t *testing.T) {
		s := Stack[float64]{}

		s.Push(1.5)
		s.Push(2.5)
		s.Push(3.5)
		require.Equal(t, 3, s.Len())

		i, err := s.Pop()
		require.Equal(t, 3.5, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, 2.5, i)
		require.Nil(t, err)

		i, err = s.Pop()
		require.Equal(t, 2.5, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, 1.5, i)
		require.Nil(t, err)

		i, err = s.Pop()
		require.Equal(t, 1.5, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, 0.0, i)
		require.Equal(t, ErrStackEmpty, err)

		i, err = s.Pop()
		require.Equal(t, 0.0, i)
		require.Equal(t, ErrStackEmpty, err)
	})
	t.Run("Int", func(t *testing.T) {
		s := Stack[int]{}

		s.Push(1)
		s.Push(2)
		s.Push(3)
		require.Equal(t, 3, s.Len())

		i, err := s.Pop()
		require.Equal(t, 3, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, 2, i)
		require.Nil(t, err)

		i, err = s.Pop()
		require.Equal(t, 2, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, 1, i)
		require.Nil(t, err)

		i, err = s.Pop()
		require.Equal(t, 1, i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, 0, i)
		require.Equal(t, ErrStackEmpty, err)

		i, err = s.Pop()
		require.Equal(t, 0, i)
		require.Equal(t, ErrStackEmpty, err)
	})
	t.Run("String", func(t *testing.T) {
		s := Stack[string]{}

		s.Push("1")
		s.Push("2")
		s.Push("3")
		require.Equal(t, 3, s.Len())

		i, err := s.Pop()
		require.Equal(t, "3", i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, "2", i)
		require.Nil(t, err)

		i, err = s.Pop()
		require.Equal(t, "2", i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, "1", i)
		require.Nil(t, err)

		i, err = s.Pop()
		require.Equal(t, "1", i)
		require.Nil(t, err)

		i, err = s.Peek()
		require.Equal(t, "", i)
		require.Equal(t, ErrStackEmpty, err)

		i, err = s.Pop()
		require.Equal(t, "", i)
		require.Equal(t, ErrStackEmpty, err)
	})
}

func TestIntStack(t *testing.T) {
	s := IntStack{}

	s.Push(1)
	s.Push(2)
	s.Push(3)
	require.Equal(t, 3, s.Len())

	i, err := s.Pop()
	require.Equal(t, 3, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, 2, i)
	require.Nil(t, err)

	i, err = s.Pop()
	require.Equal(t, 2, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, 1, i)
	require.Nil(t, err)

	i, err = s.Pop()
	require.Equal(t, 1, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, 0, i)
	require.Equal(t, ErrStackEmpty, err)

	i, err = s.Pop()
	require.Equal(t, 0, i)
	require.Equal(t, ErrStackEmpty, err)
}

func TestStringStack(t *testing.T) {
	s := StringStack{}

	s.Push("1")
	s.Push("2")
	s.Push("3")
	require.Equal(t, 3, s.Len())

	i, err := s.Pop()
	require.Equal(t, "3", i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, "2", i)
	require.Nil(t, err)

	i, err = s.Pop()
	require.Equal(t, "2", i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, "1", i)
	require.Nil(t, err)

	i, err = s.Pop()
	require.Equal(t, "1", i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, "", i)
	require.Equal(t, ErrStackEmpty, err)

	i, err = s.Pop()
	require.Equal(t, "", i)
	require.Equal(t, ErrStackEmpty, err)
}

func TestFloatStack(t *testing.T) {
	s := FloatStack{}

	s.Push(1.5)
	s.Push(2.5)
	s.Push(3.5)
	require.Equal(t, 3, s.Len())

	i, err := s.Pop()
	require.Equal(t, 3.5, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, 2.5, i)
	require.Nil(t, err)

	i, err = s.Pop()
	require.Equal(t, 2.5, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, 1.5, i)
	require.Nil(t, err)

	i, err = s.Pop()
	require.Equal(t, 1.5, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, 0.0, i)
	require.Equal(t, ErrStackEmpty, err)

	i, err = s.Pop()
	require.Equal(t, 0.0, i)
	require.Equal(t, ErrStackEmpty, err)
}

func TestBoolStack(t *testing.T) {
	s := BoolStack{}

	s.Push(true)
	s.Push(false)
	s.Push(true)
	s.Push(false)
	require.Equal(t, 4, s.Len())

	i, err := s.Pop()
	require.Equal(t, false, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, true, i)
	require.Nil(t, err)

	i, err = s.Pop()
	require.Equal(t, true, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, false, i)
	require.Nil(t, err)

	i, err = s.Pop()
	require.Equal(t, false, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, true, i)
	require.Nil(t, err)

	i, err = s.Pop()
	require.Equal(t, true, i)
	require.Nil(t, err)

	i, err = s.Peek()
	require.Equal(t, false, i)
	require.Equal(t, ErrStackEmpty, err)

	i, err = s.Pop()
	require.Equal(t, false, i)
	require.Equal(t, ErrStackEmpty, err)
}
