package queue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	t.Run("Copy", func(t *testing.T) {
		q := Queue[string]{}
		q.Enqueue("5")
		q.Enqueue("4")
		q.Enqueue("c")
		q.Enqueue("2")
		q.Enqueue("a")

		cpy := q.Copy()
		require.Equal(t, q.items, cpy.items)
	})

	t.Run("EnDeQueue", func(t *testing.T) {
		q := Queue[interface{}]{}
		input := []interface{}{"5", 4, "c", "2", "a", false, 3.14159}

		_, err := q.Dequeue()
		require.Equal(t, ErrQueueEmpty, err)

		_, err = q.Peek()
		require.Equal(t, ErrQueueEmpty, err)

		for i := 0; i < 7; i++ {
			q.Enqueue(input[i])

			itm, err := q.Peek()
			require.Nil(t, err)
			require.Equal(t, input[0], itm)
		}

		for i := 0; i < 7; i++ {
			itm, err := q.Dequeue()
			require.Nil(t, err)
			require.Equal(t, input[i], itm)
		}
	})
}
