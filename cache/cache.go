// Package cache provides a generalized cache package that allows anything to be cached
// without having to know underlying details of where that cache is stored.
// All items are serialized into JSON and stored as strings. When retrieved, it's
// unmarshaled into the provided container.
// A single cache is maintained at the package level. The provided cacher should be thread-safe,
// as no locking occurs in this package.
package cache

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/btm6084/gojson"
	"github.com/btm6084/utilities/health"
	"github.com/btm6084/utilities/metrics"
	"github.com/btm6084/utilities/redis"
)

var (
	c         Cacher
	dur       time.Duration
	deB64Hint = ":::deB64"

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

	// ForeverTTL returns the specific value that represents the no-expire TTL value.
	ForeverTTL() int
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

	switch container.(type) {
	case []byte, *[]byte:
		if !strings.HasSuffix(b, deB64Hint) {
			err = gojson.Unmarshal([]byte(b), &container)
			if err != nil {
				return err
			}

			return nil
		}

		p := reflect.ValueOf(container)
		if p.Kind() != reflect.Ptr {
			return fmt.Errorf("supplied container must be a pointer")
		}

		p = resolvePtr(p)
		if !p.CanSet() {
			return fmt.Errorf("unsettable value provided to Get")
		}

		b = strings.TrimSuffix(b, deB64Hint)

		sq := stripQuotes(b)
		if gojson.IsJSONNull(sq) {
			sq = []byte{}
		}

		d := make([]byte, base64.StdEncoding.DecodedLen(len(sq)))
		n, err := base64.StdEncoding.Decode(d, sq)

		if err != nil {
			return err
		}
		p.SetBytes(d[:n])

		return nil
	default:
		if strings.HasSuffix(b, deB64Hint) {
			b = b64DecodeString(strings.TrimSuffix(b, deB64Hint))
		}

		err = gojson.Unmarshal([]byte(b), &container)
		if err != nil {
			return err
		}

		return nil
	}
}

// Set a value in cache. Value MUST be serializeable by JSON. UnExported fields will be ignored!
func Set(m metrics.Recorder, key string, value interface{}) error {
	if !Enabled {
		return ErrCacheDisabled
	}

	return SetWithDuration(m, key, value, dur)
}

// SetWithDuration sets a value in cache. Value MUST be serializeable by JSON. UnExported fields will be ignored!
func SetWithDuration(m metrics.Recorder, key string, value interface{}, d time.Duration) error {
	if d < 0 {
		d = time.Duration(c.ForeverTTL())
	}

	if !Enabled {
		return ErrCacheDisabled
	}

	if c == nil {
		return ErrCacheNil
	}

	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	hint := ""
	switch value.(type) {
	case []byte, *[]byte:
		hint = deB64Hint
	}

	c.SetWithDuration(m, key, string(raw)+hint, d)
	return nil
}

func stripQuotes(in string) []byte {
	b := []byte(in)

	if len(b) < 2 {
		return b
	}

	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	return b
}

func b64DecodeString(in string) string {
	b := []byte(in)

	if len(b) < 2 {
		return in
	}

	reQuote := false
	if b[0] == '"' && b[len(b)-1] == '"' {
		reQuote = true
		b = b[1 : len(b)-1]
	}

	out, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return in
	}

	if reQuote {
		return `"` + string(out) + `"`
	}

	return string(out)
}

// Resolve a pointer to a concrete Value. If necessary, memory will be allocated to
// store the object being pointed to.
func resolvePtr(p reflect.Value) reflect.Value {
	op := p

	for p.Kind() == reflect.Ptr || p.Kind() == reflect.Interface {
		if p.Kind() == reflect.Ptr && !p.Elem().CanAddr() {
			child := reflect.New(p.Type().Elem()).Elem()
			p.Set(child.Addr())
		}

		if !p.Elem().IsValid() {
			break
		}

		p = p.Elem()

		// Retain the last setable value. This usually comes into play when we have
		// an interface that represents a non-settable value. The end result is we will
		// perform the extraction as if we were an interface. This is in alignment with
		// the behavior of encoding/json.Unmarshal
		if p.CanSet() {
			op = p
		}
	}

	if !p.CanSet() {
		p = op
	}

	return p
}

func HealthCheck(c Cacher) *health.Check {
	switch c.(type) {
	case *redis.Client:
		return checkRedis(c)

	case *MemoryCache:
		return checkMemoryCache(c)
	}

	return nil
}

func checkRedis(c Cacher) *health.Check {
	start := time.Now()
	hc := health.Check{
		Name:        "cache",
		Status:      health.OK,
		Description: "cache is Healthy",
		Data: map[string]interface{}{
			"cacheType": "redis",
		},
	}

	rdb, ok := c.(*redis.Client)
	if !ok {
		hc.Status = health.CRITICAL
		hc.Description = "expected redis type cacher"
		return &hc
	}

	errChan := make(chan error, 1)
	go func() {
		err := rdb.Ping(&metrics.NoOp{})
		errChan <- err
	}()

	var err error
	select {
	case err = <-errChan:
	case <-time.After(1 * time.Second):
		err = fmt.Errorf("redis cache timed out during ping operation")
	}

	if err != nil {
		hc.Status = health.CRITICAL
		hc.Description = fmt.Sprintf("redis cache storage error: %s", err)
	}

	hc.Data["pingTime"] = time.Since(start).String()
	return &hc
}

func checkMemoryCache(c Cacher) *health.Check {
	start := time.Now()
	hc := health.Check{
		Name:        "cache",
		Status:      health.OK,
		Description: "cache is Healthy",
		Data: map[string]interface{}{
			"cacheType": "memory_cache",
		},
	}

	setVal := time.Now().String()
	err := c.Set(&metrics.NoOp{}, "healthcheck_test_key", setVal)
	if err != nil {
		hc.Status = health.CRITICAL
		hc.Description = "Memory Cache Set Failed"
		return &hc
	}

	getVal, err := c.Get(&metrics.NoOp{}, "healthcheck_test_key")
	if err != nil {
		hc.Status = health.CRITICAL
		hc.Description = "Memory Cache Get Failed"
		return &hc
	}

	if g, ok := getVal.(string); ok {
		if g != setVal {
			hc.Status = health.CRITICAL
			hc.Description = "Memory Cache Get Returned Wrong Value"
			hc.Data = map[string]interface{}{
				"get": g,
				"set": setVal,
			}
			return &hc
		}
	}

	hc.Data["pingTime"] = time.Since(start).String()
	return &hc
}
