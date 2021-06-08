package logging

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTxnFields(t *testing.T) {
	ctx := ContextWithTransaction(context.Background(), "UnitTest TextTxnFields")
	txnFields := TxnFields(ctx)

	require.Equal(t, "UnitTest TextTxnFields", txnFields["txnID"])

	st := txnFields["stacktrace"].(map[string]interface{})
	require.Equal(t, st["line"], 13)

	f := st["file"].(string)
	require.True(t, strings.HasSuffix(f, "utilities/logging/transaction_test.go"), f)
}
