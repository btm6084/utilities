package metrics

// NoOp satisfies the Recorder interface, but does nothing.
type NoOp struct{}

// DatabaseSegment records a database segment that occured during the given transaction.
func (n *NoOp) DatabaseSegment(string, string, ...interface{}) func() {
	return func() {}
}

// Segment records a segment that occured during the given transaction.
func (n *NoOp) Segment(string) func() {
	return func() {}
}

// SetDBMeta no-ops.
func (n *NoOp) SetDBMeta(string, string, string) {}
