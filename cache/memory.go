package cache

import (
	"time"

	"github.com/btm6084/utilities/metrics"
	"github.com/patrickmn/go-cache"
)

// MemoryCache is an in-memory cache storage that satisfies the Cacher interface.
type MemoryCache struct {
	cache      *cache.Cache
	defaultTTL time.Duration
}

// NewMemoryCache returns an in-memory Cacher
func NewMemoryCache(defaultTTL time.Duration) Cacher {
	return &MemoryCache{defaultTTL: defaultTTL, cache: cache.New(defaultTTL, 2*defaultTTL)}
}

// ForeverTTL returns the go-cache value that represents the no-expire TTL value.
func (c *MemoryCache) ForeverTTL() int {
	return -1
}

// Get a value from cache.
func (c *MemoryCache) Get(m metrics.Recorder, key string) (interface{}, error) {
	if c.cache == nil {
		return nil, ErrCacheNil
	}

	if val, ok := c.cache.Get(key); ok {
		return val, nil
	}

	return nil, ErrNotFound

}

// Set a value in cache.
func (c *MemoryCache) Set(m metrics.Recorder, key string, value interface{}) error {
	return c.SetWithDuration(m, key, value, c.defaultTTL)
}

// SetWithDuration sets a value in cache.
func (c *MemoryCache) SetWithDuration(m metrics.Recorder, key string, value interface{}, d time.Duration) error {
	if c.cache == nil {
		return ErrCacheNil
	}

	c.cache.Set(key, value, d)
	return nil
}

// Delete removes a key from cache.
func (c *MemoryCache) Delete(m metrics.Recorder, key string) error {
	if c.cache == nil {
		return ErrCacheNil
	}

	c.cache.Delete(key)
	return nil
}
