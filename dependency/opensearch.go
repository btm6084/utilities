package dependency

import (
	"context"
	"net/http"

	opensearch "github.com/opensearch-project/opensearch-go"
)

var (
	// For now, we only allow a single search instance per service. If this changes, this can be shifted to a parameter.
	searchLabel = `searchInstance`
)

// OpensearchDependencyHandler is a middleware to inject an opensearch.Client into context.
func OpensearchDependencyHandler(search *opensearch.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ContextWithOpensearch(r.Context(), search)))
		})
	}
}

// ContextWithOpensearch returns a context containing the provided opensearch client.
func ContextWithOpensearch(ctx context.Context, search *opensearch.Client) context.Context {
	searchList, ok := ctx.Value(opensearchContextKey).(map[string]*opensearch.Client)
	if !ok || searchList == nil {
		searchList = make(map[string]*opensearch.Client)
	}

	searchList[searchLabel] = search

	nc := context.WithValue(ctx, opensearchContextKey, searchList)

	return nc
}

// OpensearchFromContext gets an opensearch client from a context.
func OpensearchFromContext(ctx context.Context) *opensearch.Client {
	searchList, ok := ctx.Value(opensearchContextKey).(map[string]*opensearch.Client)
	if !ok {
		return nil
	}

	if search, isset := searchList[searchLabel]; isset {
		return search
	}

	return nil
}
