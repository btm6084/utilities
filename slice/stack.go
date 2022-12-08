package slice

import "errors"

var (
	ErrStackEmpty = errors.New("stack is empty")
)

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(value T) {
	s.items = append(s.items, value)
}

func (s *Stack[T]) Pop() (T, error) {
	var empty T
	n := len(s.items)

	if n == 0 {
		return empty, ErrStackEmpty
	}

	i := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return i, nil
}

func (s *Stack[T]) Peek() (T, error) {
	var empty T
	n := len(s.items)

	if n == 0 {
		return empty, ErrStackEmpty
	}

	return s.items[n-1], nil
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func (s *Stack[T]) Copy() Stack[T] {
	var newStack Stack[T]
	newStack.items = make([]T, len(s.items))

	copy(newStack.items, s.items)
	return newStack
}

type IntStack struct {
	ints []int
}

func (s *IntStack) Push(i int) {
	s.ints = append(s.ints, i)
}

func (s *IntStack) Pop() (int, error) {
	if len(s.ints) == 0 {
		return 0, ErrStackEmpty
	}

	i := s.ints[len(s.ints)-1]
	s.ints = s.ints[:len(s.ints)-1]
	return i, nil
}

func (s *IntStack) Peek() (int, error) {
	if len(s.ints) == 0 {
		return 0, ErrStackEmpty
	}

	return s.ints[len(s.ints)-1], nil
}

func (s *IntStack) Len() int {
	return len(s.ints)
}

type StringStack struct {
	strings []string
}

func (s *StringStack) Push(i string) {
	s.strings = append(s.strings, i)
}

func (s *StringStack) Pop() (string, error) {
	if len(s.strings) == 0 {
		return "", ErrStackEmpty
	}

	i := s.strings[len(s.strings)-1]
	s.strings = s.strings[:len(s.strings)-1]
	return i, nil
}

func (s *StringStack) Peek() (string, error) {
	if len(s.strings) == 0 {
		return "", ErrStackEmpty
	}

	return s.strings[len(s.strings)-1], nil
}

func (s *StringStack) Len() int {
	return len(s.strings)
}

type FloatStack struct {
	floats []float64
}

func (s *FloatStack) Push(i float64) {
	s.floats = append(s.floats, i)
}

func (s *FloatStack) Pop() (float64, error) {
	if len(s.floats) == 0 {
		return 0, ErrStackEmpty
	}

	i := s.floats[len(s.floats)-1]
	s.floats = s.floats[:len(s.floats)-1]
	return i, nil
}

func (s *FloatStack) Peek() (float64, error) {
	if len(s.floats) == 0 {
		return 0, ErrStackEmpty
	}

	return s.floats[len(s.floats)-1], nil
}

func (s *FloatStack) Len() int {
	return len(s.floats)
}

type BoolStack struct {
	bools []bool
}

func (s *BoolStack) Push(i bool) {
	s.bools = append(s.bools, i)
}

func (s *BoolStack) Pop() (bool, error) {
	if len(s.bools) == 0 {
		return false, ErrStackEmpty
	}

	i := s.bools[len(s.bools)-1]
	s.bools = s.bools[:len(s.bools)-1]
	return i, nil
}

func (s *BoolStack) Peek() (bool, error) {
	if len(s.bools) == 0 {
		return false, ErrStackEmpty
	}

	return s.bools[len(s.bools)-1], nil
}

func (s *BoolStack) Len() int {
	return len(s.bools)
}
