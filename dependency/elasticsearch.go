package dependency

import (
	"context"
	"net/http"

	"github.com/elastic/go-elasticsearch"
)

// ElasticsearchDependencyHandler is a middleware to inject an elasticsearch.Client into context.
func ElasticsearchDependencyHandler(es *elasticsearch.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ContextWithElasticsearch(r.Context(), es)))
		})
	}
}

// ConextWithElasticsearch returns a context containing the provided elasticsearch client.
func ContextWithElasticsearch(ctx context.Context, es *elasticsearch.Client) context.Context {
	label := `esInstance`

	esList, ok := ctx.Value(esContextKey).(map[string]*elasticsearch.Client)
	if !ok || esList == nil {
		esList = make(map[string]*elasticsearch.Client)
	}

	esList[label] = es

	nc := context.WithValue(ctx, esContextKey, esList)

	return nc
}

func ESFromContext(ctx context.Context, label string) *elasticsearch.Client {
	esList, ok := ctx.Value(esContextKey).(map[string]*elasticsearch.Client)
	if !ok {
		return nil
	}

	if es, isset := esList[label]; isset {
		return es
	}

	return nil
}
