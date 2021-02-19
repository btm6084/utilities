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

	// ErrCacheDisabled is returned when caching falied due to caching being disabled
	ErrCacheDisabled = errors.New("cache is disabled")

	// ErrDeserialize is returned when cache extraction failed due to deserialization errors.
	ErrDeserialize = errors.New("cached value could not be deserialized")

	// ErrNotFound is returned when no data was found
	ErrNotFound = errors.New("not found")

	// Enabled allows caching to be turned on and off. This is useful for turning cache off
	// via environment variables.
	Enabled = true
)

// Cacher provides an interface for working with a cache store.
type Cacher interface {
	Get(metrics.Recorder, string) (interface{}, error)
	Set(metrics.Recorder, string, interface{}) error
	SetWithDuration(metrics.Recorder, string, interface{}, time.Duration) error
	Delete(metrics.Recorder, string) error
}

func init() {
	if !Enabled {
		return
	}

	// Start with an initial memory cache. This can be tailored by calling Initialize.
	c = NewMemoryCache(5 * time.Minute)
}

// Initialize must be called prior to use. Do this in main.
func Initialize(cache Cacher, defaultDuration time.Duration) {
	if !Enabled {
		return
	}

	c = cache
	dur = defaultDuration
}

// Get a value from cache.
func Get(r metrics.Recorder, key string, container interface{}) error {
	if !Enabled {
		return ErrCacheDisabled
	}

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
	if !Enabled {
		return ErrCacheDisabled
	}

	return SetWithDuration(m, key, value, dur)
}

// SetWithDuration sets a value in cache.
func SetWithDuration(m metrics.Recorder, key string, value interface{}, d time.Duration) error {
	if !Enabled {
		return ErrCacheDisabled
	}

	if c == nil {
		return ErrCacheNil
	}

	v := value
	if d, ok := value.([]byte); ok {
		v = string(d)
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.SetWithDuration(m, key, string(raw), d)
	return nil
}
