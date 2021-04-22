package dependency

type key int

const (
	// dbContextKey is the key under which we store the database interface in context.
	dbContextKey key = 0

	// esContextKey is the key under which we store the elasticsearch interface in context.
	esContextKey key = 1
)
