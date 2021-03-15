package dependency

type key int

const (
	// dbContextKey is the key under which we store the database interface in context.
	dbContextKey key = 0
)
