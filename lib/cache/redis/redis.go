// Package redis provides a Redis-backed implementation of CachingService.
// This is not currently used. But may be implemented in the future if we decide that
// database storage is too durable for certain types of caching activities.
package redis

import (
	"errors"
	"time"

	redis "github.com/gomodule/redigo/redis"

	redigotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/garyburd/redigo"
)

var (
	errNotFound = errors.New("not found")
)

// Option is a function that modifies a redis.Pool to override default values.
type Option func(p *redis.Pool)

// WithMaxIdle overrides the default maximum number of idle connections.
func WithMaxIdle(n int) Option {
	return func(p *redis.Pool) {
		p.MaxIdle = n
	}
}

// WithTimeout overrides the default timeout for operations.
func WithTimeout(d time.Duration) Option {
	return func(p *redis.Pool) {
		p.IdleTimeout = d
	}
}

// DefaultMaxIdleConnections controls how many idle Redis connections we allow.
const DefaultMaxIdleConnections = 10

// DefaultTimeout controls how long we wait for operations to succeed before
// erroring.
const DefaultTimeout = 5 * time.Second

func NewPool(url string, opts ...Option) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     DefaultMaxIdleConnections,
		IdleTimeout: DefaultTimeout,
		Dial: func() (redis.Conn, error) {
			return redigotrace.DialURL(url)
		},
	}
	for _, opt := range opts {
		opt(pool)
	}
	return pool
}
