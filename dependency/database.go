package dependency

import (
	"context"
	"net/http"

	"github.com/btm6084/godb"
)

// DatabaseDependencyHandler is a middleware to inject a database into context.
func DatabaseDependencyHandler(label string, db godb.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ContextWithDatabase(r.Context(), label, db)))
		})
	}
}

// ConextWithDatabase returns a context containing the provided database.
func ContextWithDatabase(ctx context.Context, label string, db godb.Database) context.Context {
	dbList, ok := ctx.Value(dbContextKey).(map[string]godb.Database)
	if !ok || dbList == nil {
		dbList = make(map[string]godb.Database)
	}

	dbList[label] = db

	nc := context.WithValue(ctx, dbContextKey, dbList)

	return nc
}

func DBFromContext(ctx context.Context, label string) godb.Database {
	dbList, ok := ctx.Value(dbContextKey).(map[string]godb.Database)
	if !ok {
		return nil
	}

	if db, isset := dbList[label]; isset {
		return db
	}

	return nil
}
