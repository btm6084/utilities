package metrics

import (
	"context"
)

// Recorder is an interface for recording metrics about the application.
type Recorder interface {
	// SetDBMeta assigns the Database meta data for the next DatabaseSegment call. Call before calling DatabaseSegment.
	SetDBMeta(string, string, string)

	// DatabaseSegment records a database segment that occured during the given transaction.
	DatabaseSegment(string, string, ...interface{}) func()

	// Segment records a segment that occured during the given transaction.
	Segment(string) func()
}

var (
	// MetricsRecorder is used by the GetRecorder function to determine which recorder
	// to return. This allows an application to set a default recorder and simply call
	// GetRecorder.
	// Enum:
	//    noop
	//    newrelic
	MetricsRecorder = "noop"
)

// GetRecorder returns an appropriate recorder based on the value of MetricsRecorder.
// Populate MetricsRecorder during setup in main to change which recorder is returned.
func GetRecorder(ctx context.Context) Recorder {
	switch MetricsRecorder {
	case "newrelic":
		return NewRelicFromContext(ctx)
	default:
		return &NoOp{}
	}
}
