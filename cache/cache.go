package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	c   *cache.Cache
	dur time.Duration
)

func init() {
	dur = time.Duration(300 * time.Second)
	c = cache.New(dur, 2*dur)
}

// Get a value from cache.
func Get(key string) (interface{}, bool) {
	return c.Get(key)
}

// Set a value in cache.
func Set(key string, value interface{}) {
	c.SetDefault(key, value)
}

// SetWithDuration sets a value in cache.
func SetWithDuration(key string, value interface{}, d time.Duration) {
	c.Set(key, value, d)
}
