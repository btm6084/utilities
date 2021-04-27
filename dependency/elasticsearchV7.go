package dependency

import (
	"context"
	"net/http"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

var (
	// For now, we only allow a single es instance per service. If this changes, this can be shifted to a parameter.
	esV7Label = `esV7Instance`
)

// ElasticsearchV7DependencyHandler is a middleware to inject an elasticsearchv7.Client into context.
func ElasticsearchV7DependencyHandler(es *elasticsearch7.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ContextWithElasticsearchV7(r.Context(), es)))
		})
	}
}

// ContextWithElasticsearchV7 returns a context containing the provided elasticsearch client.
func ContextWithElasticsearchV7(ctx context.Context, es *elasticsearch7.Client) context.Context {
	esList, ok := ctx.Value(esV7ContextKey).(map[string]*elasticsearch7.Client)
	if !ok || esList == nil {
		esList = make(map[string]*elasticsearch7.Client)
	}

	esList[esV7Label] = es

	nc := context.WithValue(ctx, esV7ContextKey, esList)

	return nc
}

// ESV7FromContext gets an elasticsearch client from a context.
func ESV7FromContext(ctx context.Context) *elasticsearch7.Client {
	esList, ok := ctx.Value(esV7ContextKey).(map[string]*elasticsearch7.Client)
	if !ok {
		return nil
	}

	if es, isset := esList[esV7Label]; isset {
		return es
	}

	return nil
}
