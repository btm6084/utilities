package logging

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type key int

const (
	// txnIDKey is the key under which we store the transaction ID in context.
	txnIDKey key = 0
)

// TransactionHandler creates a unique ID for every request.
func TransactionHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			txnID := uuid.New().String()

			c := req.Context()
			nc := context.WithValue(c, txnIDKey, txnID)
			next.ServeHTTP(w, req.WithContext(nc))
		})
	}
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
func TxnFields(r *http.Request) log.Fields {
	return log.Fields{"txnID": TransactionFromContext(r.Context())}
}
