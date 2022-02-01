package redis

import (
	"time"

	"github.com/btm6084/utilities/metrics"
)

type Cache interface {
	ForeverTTL() int
	Ping(metrics.Recorder) error
	Get(metrics.Recorder, string) (interface{}, error)
	TTL(metrics.Recorder, string) (time.Duration, error)
	Set(metrics.Recorder, string, interface{}) error
	SetWithDuration(metrics.Recorder, string, interface{}, time.Duration) error
	Delete(metrics.Recorder, string) error
	IncrementHash(metrics.Recorder, string, string, int) error
	IncrementHashWithDuration(metrics.Recorder, string, string, int, time.Duration) error
	GetHashSet(metrics.Recorder, []string) ([]map[string]string, error)
}
