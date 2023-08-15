package logging

import (
	"context"
	"net/http"

	"github.com/btm6084/utilities/stack"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type key int

const (
	// txnIDKey is the key under which we store the transaction ID in context.
	txnIDKey key = 0

	// pageRequestIDKey is the key under which we store the per-page request ID in context.
	pageRequestIDKey key = 0
)

// TransactionHandler creates a unique ID for every request.
func AllTransactionHandler(txnIDHeader, pageReqIDHeader string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			txnID := req.Header.Get(txnIDHeader)
			if txnID == "" {
				txnID = uuid.New().String()
			}

			pageReqID := req.Header.Get(pageReqIDHeader)
			if pageReqID == "" {
				pageReqID = uuid.New().String()
			}

			ctx := ContextWithTransaction(req.Context(), txnID)
			ctx = ContextWithPageRequestID(ctx, pageReqID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// TransactionHandler creates a unique ID for every request.
func TransactionHandler(header string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			txnID := req.Header.Get(header)
			if txnID == "" {
				txnID = uuid.New().String()
			}

			next.ServeHTTP(w, req.WithContext(ContextWithTransaction(req.Context(), txnID)))
		})
	}
}

func ContextWithTransaction(ctx context.Context, txnID string) context.Context {
	return context.WithValue(ctx, txnIDKey, txnID)
}

// TransactionFromContext retrieves the txnID from the given context.
func TransactionFromContext(ctx context.Context) string {
	txnID, ok := ctx.Value(txnIDKey).(string)

	if !ok {
		return ""
	}

	return txnID
}

// PageRequestIDHandler creates a unique ID for every request.
func PageRequestIDHandler(header string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			pageReqID := req.Header.Get(header)
			if pageReqID == "" {
				pageReqID = uuid.New().String()
			}

			next.ServeHTTP(w, req.WithContext(ContextWithPageRequestID(req.Context(), pageReqID)))
		})
	}
}

func ContextWithPageRequestID(ctx context.Context, pageReqID string) context.Context {
	return context.WithValue(ctx, pageRequestIDKey, pageReqID)
}

// PageRequestIDFromContext retrieves the pageReqID from the given context.
func PageRequestIDFromContext(ctx context.Context) string {
	pageReqID, ok := ctx.Value(pageRequestIDKey).(string)

	if !ok {
		return ""
	}

	return pageReqID
}

// TxnFields builds a log.Fields object when all you need is a transactionID.
func TxnFields(ctx context.Context) log.Fields {
	f, l := stack.Trace(1) // 1 to refer to the caller of this fn.
	return log.Fields{"txnID": TransactionFromContext(ctx), "pageReqID": PageRequestIDFromContext(ctx), "stacktrace": map[string]interface{}{"file": f, "line": l}}
}
