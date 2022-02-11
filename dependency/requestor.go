package dependency

import (
	"context"
	"net/http"

	httputil "github.com/btm6084/utilities/http"
)

// RequestorDependencyHandler is a middleware to inject a requestor into context.
func RequestorDependencyHandler(label string, req httputil.Requestor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(ContextWithRequestor(r.Context(), label, req)))
		})
	}
}

// ContextWithRequestor returns a context containing the provided requestor.
func ContextWithRequestor(ctx context.Context, label string, req httputil.Requestor) context.Context {
	reqList, ok := ctx.Value(reqContextKey).(map[string]httputil.Requestor)
	if !ok || reqList == nil {
		reqList = make(map[string]httputil.Requestor)
	}

	reqList[label] = req

	nc := context.WithValue(ctx, reqContextKey, reqList)

	return nc
}

// RequestorFromContext returns a requestor from a context.
func RequestorFromContext(ctx context.Context, label string) httputil.Requestor {
	reqList, ok := ctx.Value(reqContextKey).(map[string]httputil.Requestor)
	if !ok {
		return nil
	}

	if req, isset := reqList[label]; isset {
		return req
	}

	return nil
}
