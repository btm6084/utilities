package logging

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxnFields(t *testing.T) {
	ctx := ContextWithTransaction(context.Background(), "UnitTest TextTxnFields")
	txnFields := TxnFields(ctx)

	assert.Equal(t, "UnitTest TextTxnFields", txnFields["txnID"])

	st := txnFields["stacktrace"].(map[string]interface{})
	assert.Equal(t, st["line"], 13)

	f := st["file"].(string)
	assert.True(t, strings.HasSuffix(f, "utilities/logging/transaction_test.go"), f)
}
