// Package redis is a facade in front of github.com/go-redis/redis/v8
package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/btm6084/utilities/metrics"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
)

var (
	// DefaultTTL defines a default TTL value for all keys.
	DefaultTTL time.Duration

	// Namespace allows for a custom Namespace to be added to all keys.
	Namespace = ""

	// ErrNotFound is returned when no data was found
	ErrNotFound = errors.New("not found")

	// Compiler will enforce the interface and let us know if the contract is broken.
	_ Cache = (*Client)(nil)
)

func init() {
	// Setup an 5m initial TTL. This can be overwritten, or ignored by using SetWithDuration.
	DefaultTTL = 5 * time.Minute
}

// Client provides interaction with redis.
type Client struct {
	RDB            *redis.Client
	requestTimeout time.Duration
}

// New creates a new client.
func New(server string, requestTimeout time.Duration, clientName string) *Client {
	return &Client{
		RDB: redis.NewClient(&redis.Options{
			Addr:      server,
			Password:  "", // no password set
			DB:        0,  // use default DB
			OnConnect: onConnect(clientName),
		}),
		requestTimeout: requestTimeout,
	}
}

func onConnect(clientName string) func(context.Context, *redis.Conn) error {
	return func(ctx context.Context, cn *redis.Conn) error {
		p := cn.Pipeline()
		csn := p.ClientSetName(ctx, clientName)
		p.Exec(ctx)

		return csn.Err()
	}
}

// ForeverTTL returns the redis value that represents the no-expire TTL value.
func (c *Client) ForeverTTL() int {
	return 0
}

// Ping tests the connection to Redis.
func (c *Client) Ping(r metrics.Recorder) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	r.SetDBMeta("Redis", "ping", "PING")
	defer r.DatabaseSegment("redis", "ping")()
	sts := c.RDB.Ping(ctx)
	if sts.Err() != nil {
		return sts.Err()
	}

	if sts.Val() != "PONG" {
		return errors.New("unable to ping redis")
	}

	return nil
}

func (c *Client) GetString(r metrics.Recorder, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	key = Namespace + key

	r.SetDBMeta("Redis", key, "GET")
	defer r.DatabaseSegment("redis", "get key")()
	rsp := c.RDB.Get(ctx, key)
	if rsp.Err() != nil {
		if rsp.Err() == redis.Nil {
			return "", ErrNotFound
		}
		return "", rsp.Err()
	}

	return rsp.Val(), nil
}

// Get returns the value at `key`
func (c *Client) Get(r metrics.Recorder, key string) (interface{}, error) {
	return c.GetString(r, key)
}

// TTL returns the TTL for the value at `key`
func (c *Client) TTL(r metrics.Recorder, key string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	key = Namespace + key

	r.SetDBMeta("Redis", key, "TTL")
	defer r.DatabaseSegment("redis", "get ttl")()
	rsp := c.RDB.TTL(ctx, key)
	if rsp.Err() != nil {
		return 1 * time.Microsecond, rsp.Err()
	}

	return rsp.Val(), nil
}

// Set stores the value at `key` with the default TTL.
func (c *Client) Set(r metrics.Recorder, key string, value interface{}) error {
	return c.SetWithDuration(r, key, value, DefaultTTL)
}

// SetWithDuration stores the value at `key` with the provided TTL.
func (c *Client) SetWithDuration(r metrics.Recorder, key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	key = Namespace + key

	r.SetDBMeta("Redis", key, "SET")
	defer r.DatabaseSegment("redis", "set with duration", value, ttl)()
	rsp := c.RDB.Set(ctx, key, value, ttl)
	if rsp.Err() != nil && rsp.Err() != redis.Nil {
		return rsp.Err()
	}

	return nil
}

// Delete removes the value at `key`
func (c *Client) Delete(r metrics.Recorder, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	key = Namespace + key

	r.SetDBMeta("Redis", key, "DEL")
	defer r.DatabaseSegment("redis", "del key")()
	rsp := c.RDB.Del(ctx, key)
	if rsp.Err() != nil && rsp.Err() != redis.Nil {
		return rsp.Err()
	}

	return nil
}

// IncrementHash increments a value at a give key/field location. Hash will live for the standard ttl duration.
func (c *Client) IncrementHash(r metrics.Recorder, key, field string, amount int) error {
	return c.IncrementHashWithDuration(r, key, field, amount, DefaultTTL)
}

// IncrementHashWithDuration increments a value at a give key/field location.
//
// We can only expire the entire hash. We set a TTL via Expire once only if the hash key has not
// previously existed. We make the assumption that if it already exists, it's already had an Expire set.
func (c *Client) IncrementHashWithDuration(r metrics.Recorder, key, field string, amount int, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	key = Namespace + key

	expire := false
	end := r.DatabaseSegment("redis", "exists", key)
	if c.RDB.Exists(ctx, key).Val() == 0 {
		expire = true
	}
	end()

	r.SetDBMeta("Redis", key, "HINCRBY")
	end = r.DatabaseSegment("redis", "hash increment at key.field", field, amount)
	rsp := c.RDB.HIncrBy(ctx, key, field, int64(amount))
	end()
	if rsp.Err() != nil && rsp.Err() != redis.Nil {
		return rsp.Err()
	}

	if expire {
		end = r.DatabaseSegment("redis", "expire", key)
		c.RDB.Expire(ctx, key, ttl)
		end()
	}

	return nil
}

// GetHashSet uses a redis pipe to retrieve a number of hashes and return an aggregate of their responses.
func (c *Client) GetHashSet(r metrics.Recorder, keys []string) ([]map[string]string, error) {
	if len(keys) == 0 {
		return nil, ErrNotFound
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	r.SetDBMeta("Redis", strings.Join(keys, ","), "HGETALL PIPE "+cast.ToString(len(keys)))
	defer r.DatabaseSegment("redis", "get hash set")()

	pipe := c.RDB.TxPipeline()

	var vals []*redis.StringStringMapCmd

	for _, key := range keys {
		key = Namespace + key
		vals = append(vals, pipe.HGetAll(ctx, key))
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var result []map[string]string

	for _, v := range vals {
		result = append(result, v.Val())
	}

	return result, nil
}
