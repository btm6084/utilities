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
)

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

// TxnFields builds a log.Fields object when all you need is a transactionID.
func TxnFields(ctx context.Context) log.Fields {
	f, l := stack.Trace(1) // 1 to refer to the caller of this fn.
	return log.Fields{"txnID": TransactionFromContext(ctx), "stacktrace": map[string]interface{}{"file": f, "line": l}}
}
