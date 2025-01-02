// Package cache provides a simple interface for implementing caching in applications
// By default, an in-memory backed cache is created at startup
package cache

import (
	"context"
	"time"

	"goeasy.dev/cache/providers/localcache"
	"goeasy.dev/errors"
	"goeasy.dev/observability/metrics"
)

func init() {
	var err error
	defaultProvider, err = localcache.NewLocal(context.Background(), time.Minute*5)
	if err != nil {
		panic(err)
	}
}

var (
	defaultProvider   CacheProvider
	disableCacheKey   = struct{}{}
	cachePutMetric    = metrics.NewCounter("cache_put")
	cachePutErrMetric = metrics.NewCounter("cache_put_error")
	cacheHitMetric    = metrics.NewCounter("cache_hit")
	cacheMissMetric   = metrics.NewCounter("cache_miss")
	cacheErrMetric    = metrics.NewCounter("cache_error")
)

// Keyed returns an ID that can be used as the cache key
type Keyed interface {
	ID() string
}

// CacheProvider is implemented by backing cache stores such as Redis
type CacheProvider interface {
	Put(context.Context, string, interface{}) error
	Get(context.Context, string, interface{}) error
}

// Put places a value in the cache.
// Put can be invoked by either supplying the key and value, or by supplying a value that implements `Keyed`
//
//	cache.Put(ctx, "my-key", "my-value")
//	cache.Put(ctx, valueThatImplementsKeyed)
func Put(ctx context.Context, args ...interface{}) error {
	if isCacheDisabled(ctx) {
		return nil
	}

	key, value, err := parseArgs(args)
	if err != nil {
		return err
	}

	err = defaultProvider.Put(ctx, key, value)
	if err != nil {
		cachePutErrMetric.Inc(ctx)
		return errors.Wrap(err, "unable to store value with provider")
	}

	cachePutMetric.Inc(ctx)
	return nil
}

// Get retrives a value from the cache.
// Get can be invoked by either supplying the key and value, or by supplying a value that implements `Keyed`
// When using the later invokation, the ID() method needs to return the id that is the key of the object
//
//	cache.Get(ctx, "my-key", &dest)
//	cache.Put(ctx, &valueThatImplementsKeyed)
//
// If the value is not found in the cache, errors.ErrNotFound will be returned, nil otherwise
func Get(ctx context.Context, args ...interface{}) error {
	if isCacheDisabled(ctx) {
		return errors.ErrNotFound
	}

	key, dest, err := parseArgs(args)
	if err != nil {
		return err
	}

	err = defaultProvider.Get(ctx, key, dest)
	if errors.Is(err, errors.ErrNotFound) {
		cacheMissMetric.Inc(ctx)
		return err
	} else if err != nil {
		cacheErrMetric.Inc(ctx)
		return errors.Wrap(err, "unable to retrieve data from provider")
	}

	cacheHitMetric.Inc(ctx)
	return nil
}

// SetProvider changes the CacheProvider that is used for cache operations
func SetProvider(provider CacheProvider) {
	defaultProvider = provider
}

// DisableCacheFor will cause all cache operations using the context to be a no-op
func DisableCacheFor(ctx context.Context) context.Context {
	return context.WithValue(ctx, disableCacheKey, true)
}

func parseArgs(args []interface{}) (string, interface{}, error) {
	if len(args) == 2 {
		key, ok := args[0].(string)
		if !ok {
			return "", nil, errors.New("when supplying multiple arguements, first arguement must be a string")
		}

		return key, args[1], nil
	}

	keyed, ok := args[0].(Keyed)
	if !ok {
		return "", nil, errors.New("when supplying a single arguement, it must implement cache.Keyed")
	}

	return keyed.ID(), args[0], nil
}

func isCacheDisabled(ctx context.Context) bool {
	disabled, ok := ctx.Value(disableCacheKey).(bool)
	if !ok {
		return false
	}

	return disabled
}
