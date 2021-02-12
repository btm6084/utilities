package metrics

import (
	"context"

	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/cast"
)

// NewRelic implements the transaction.metrics interface.
// Populate DBName, Table, Operation before calling DatabaseSegment.
type NewRelic struct {
	Transaction *newrelic.Transaction
	DBName      string
	Collection  string
	Operation   string
}

// NewRelicFromContext creates a new NewRelic recorder from the given context, or a NoOp recorder
// if there is no NewRelic transaction in the context.
func NewRelicFromContext(ctx context.Context) Recorder {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return &NoOp{}
	}

	return &NewRelic{Transaction: txn}
}

// SetDBMeta assigns the Database meta data for the next DatabaseSegment call. Call before calling DatabaseSegment.
func (nr *NewRelic) SetDBMeta(db, collection, operation string) {
	nr.DBName = db
	nr.Collection = collection
	nr.Operation = operation
}

// DatabaseSegment records a database segment that occured during the given transaction.
func (nr *NewRelic) DatabaseSegment(product, query string, args ...interface{}) func() {
	var pdt newrelic.DatastoreProduct

	switch product {
	case "mssql":
		pdt = newrelic.DatastoreMSSQL
	case "postgres":
		pdt = newrelic.DatastorePostgres
	case "redis":
		pdt = newrelic.DatastoreRedis
	}

	s := newrelic.DatastoreSegment{
		StartTime:          newrelic.StartSegmentNow(nr.Transaction),
		Product:            pdt,
		Collection:         nr.Collection,
		Operation:          nr.Operation,
		ParameterizedQuery: query,
		QueryParameters:    argsToMap(args),
		DatabaseName:       nr.DBName,
	}
	return s.End
}

// Segment records a segment that occured during the given transaction.
func (nr *NewRelic) Segment(name string) func() {
	return newrelic.StartSegment(nr.Transaction, name).End
}

func argsToMap(args []interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range args {
		m["p"+cast.ToString(k)] = v
	}

	return m
}
