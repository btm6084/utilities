package dependency

import (
	"context"
	"net/http"

	elasticsearch6 "github.com/elastic/go-elasticsearch/v6"
)

var (
	// For now, we only allow a single es instance per service. If this changes, this can be shifted to a parameter.
	esV6Label = `esV6Instance`
)

// ElasticsearchV6DependencyHandler is a middleware to inject an elasticsearchv6.Client into context.
func ElasticsearchV6DependencyHandler(es *elasticsearch6.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ContextWithElasticsearchV6(r.Context(), es)))
		})
	}
}

// ContextWithElasticsearchV6 returns a context containing the provided elasticsearch client.
func ContextWithElasticsearchV6(ctx context.Context, es *elasticsearch6.Client) context.Context {
	esList, ok := ctx.Value(esV6ContextKey).(map[string]*elasticsearch6.Client)
	if !ok || esList == nil {
		esList = make(map[string]*elasticsearch6.Client)
	}

	esList[esV6Label] = es

	nc := context.WithValue(ctx, esV6ContextKey, esList)

	return nc
}

// ESV6FromContext gets an elasticsearch client from a context.
func ESV6FromContext(ctx context.Context) *elasticsearch6.Client {
	esList, ok := ctx.Value(esV6ContextKey).(map[string]*elasticsearch6.Client)
	if !ok {
		return nil
	}

	if es, isset := esList[esV6Label]; isset {
		return es
	}

	return nil
}
