// Package logging > cache provides information about the cached nature of the request.
// This exists in package logging instead cache due to circular dependency problems
// coupled with the fact that we're trying to inform the logger with this information
package logging

import (
	"context"
	"net/http"
)

const (
	// cacheStatusKey is the key under which we store the cache status in context.
	cacheStatusKey key = 0
)

// CacheStatusFromContext retrieves the cache status from the given context.
func CacheStatusFromContext(ctx context.Context) bool {
	cs, ok := ctx.Value(cacheStatusKey).(bool)

	if !ok {
		return false
	}

	return cs
}

// RequestWithCacheStatus attaches a cache status to the request.
func RequestWithCacheStatus(r *http.Request, status bool) *http.Request {
	c := r.Context()
	nc := context.WithValue(c, cacheStatusKey, status)

	return r.WithContext(nc)
}
