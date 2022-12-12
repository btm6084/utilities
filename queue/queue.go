package queue

import "errors"

var (
	ErrQueueEmpty = errors.New("queue is empty")
)

type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Enqueue(value T) {
	q.items = append(q.items, value)
}

func (q *Queue[T]) Dequeue() (T, error) {
	var empty T
	if len(q.items) == 0 {
		return empty, ErrQueueEmpty
	}

	i := q.items[0]
	q.items = q.items[1:]
	return i, nil
}

func (q *Queue[T]) Peek() (T, error) {
	var empty T

	if len(q.items) == 0 {
		return empty, ErrQueueEmpty
	}

	return q.items[0], nil
}

func (q *Queue[T]) Len() int {
	return len(q.items)
}

func (q *Queue[T]) Copy() Queue[T] {
	var newQueue Queue[T]
	newQueue.items = make([]T, len(q.items))

	copy(newQueue.items, q.items)
	return newQueue
}
