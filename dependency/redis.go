package dependency

import (
	"context"
	"net/http"

	"github.com/btm6084/utilities/redis"
)

var (
	// For now, we only allow a single RDB instance per service. If this changes, this can be shifted to a parameter.
	rdbLabel = `rdbInstance`
)

// RedisDependencyHandler is a middleware to inject a redis.Client into context.
func RedisDependencyHandler(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ContextWithRedis(r.Context(), rdb)))
		})
	}
}

// ConextWithRedis returns a context containing the provided redis client.
func ContextWithRedis(ctx context.Context, rdb *redis.Client) context.Context {
	rdbList, ok := ctx.Value(rdbContextKey).(map[string]*redis.Client)
	if !ok || rdbList == nil {
		rdbList = make(map[string]*redis.Client)
	}

	rdbList[rdbLabel] = rdb

	nc := context.WithValue(ctx, rdbContextKey, rdbList)

	return nc
}

func RDBFromContext(ctx context.Context) *redis.Client {
	rdbList, ok := ctx.Value(rdbContextKey).(map[string]*redis.Client)
	if !ok {
		return nil
	}

	if rdb, isset := rdbList[rdbLabel]; isset {
		return rdb
	}

	return nil
}
