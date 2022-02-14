package redis

import (
	"time"

	"github.com/btm6084/utilities/metrics"
)

var (
	// Compiler will enforce the interface and let us know if the contract is broken.
	_ Cache = (*Noop)(nil)
)

// Noop allows us to have a passthrough, do nothing cache.
type Noop struct{}

func (n *Noop) ForeverTTL() int                                           { return 0 }
func (n *Noop) Ping(metrics.Recorder) error                               { return nil }
func (n *Noop) GetString(metrics.Recorder, string) (string, error)        { return "", ErrNotFound }
func (n *Noop) Get(metrics.Recorder, string) (interface{}, error)         { return nil, ErrNotFound }
func (n *Noop) TTL(metrics.Recorder, string) (time.Duration, error)       { return 0, nil }
func (n *Noop) Set(metrics.Recorder, string, interface{}) error           { return nil }
func (n *Noop) Delete(metrics.Recorder, string) error                     { return nil }
func (n *Noop) IncrementHash(metrics.Recorder, string, string, int) error { return nil }
func (n *Noop) IncrementHashWithDuration(metrics.Recorder, string, string, int, time.Duration) error {
	return nil
}
func (n *Noop) GetHashSet(metrics.Recorder, []string) ([]map[string]string, error) {
	return nil, ErrNotFound
}
func (n *Noop) SetWithDuration(metrics.Recorder, string, interface{}, time.Duration) error {
	return nil
}
