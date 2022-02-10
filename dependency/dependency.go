package dependency

type key int

const (
	// dbContextKey is the key under which we store the database interface in context.
	dbContextKey key = 0

	// esV6ContextKey is the key under which we store the elasticsearch v6 interface in context.
	esV6ContextKey key = 1

	// esV7ContextKey is the key under which we store the elasticsearch v7 interface in context.
	esV7ContextKey key = 2

	// rdbContextKey is the key under which we store the redis interface in context.
	rdbContextKey key = 3

	// reqContextKey is the key under which we store the requestor interface in context.
	reqContextKey key = 4
)
