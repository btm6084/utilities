// Package cache provides a generalized cache package that allows anything to be cached
// without having to know underlying details of where that cache is stored.
// All items are serialized into JSON and stored as strings. When retrieved, it's
// unmarshaled into the provided container.
// A single cache is maintained at the package level. The provided cacher should be thread-safe,
// as no locking occurs in this package.
package cache

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/btm6084/gojson"
	"github.com/btm6084/utilities/metrics"
)

var (
	c   Cacher
	dur time.Duration

	// ErrCacheNil is returned when caching falied due to the provided Cacher being nil
	ErrCacheNil = errors.New("cache is nil")

	// ErrDeserialize is returned when cache extraction failed due to deserialization errors.
	ErrDeserialize = errors.New("cached value could not be deserialized")

	// ErrNotFound is returned when no data was found
	ErrNotFound = errors.New("not found")
)

// Cacher provides an interface for working with a cache store.
type Cacher interface {
	Get(metrics.Recorder, string) (interface{}, error)
	Set(metrics.Recorder, string, interface{}) error
	SetWithDuration(metrics.Recorder, string, interface{}, time.Duration) error
	Delete(metrics.Recorder, string) error
}

func init() {
	// Start with an initial memory cache. This can be tailored by calling Initialize.
	c = NewMemoryCache(5 * time.Minute)
}

// Initialize must be called prior to use. Do this in main.
func Initialize(cache Cacher, defaultDuration time.Duration) {
	c = cache
	dur = defaultDuration
}

// Get a value from cache.
func Get(r metrics.Recorder, key string, container interface{}) error {
	if c == nil {
		return ErrCacheNil
	}

	raw, err := c.Get(r, key)
	if err != nil {
		return err
	}

	b, ok := raw.(string)
	if !ok {
		return ErrDeserialize
	}

	err = gojson.Unmarshal([]byte(b), &container)
	if err != nil {
		return err
	}

	return nil
}

// Set a value in cache.
func Set(m metrics.Recorder, key string, value interface{}) error {
	return SetWithDuration(m, key, value, dur)
}

// SetWithDuration sets a value in cache.
func SetWithDuration(m metrics.Recorder, key string, value interface{}, d time.Duration) error {
	if c == nil {
		return ErrCacheNil
	}

	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.SetWithDuration(m, key, string(v), d)
	return nil
}
