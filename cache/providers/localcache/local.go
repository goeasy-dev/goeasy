package localcache

import (
	"context"
	"time"

	"goeasy.dev/errors"
	"goeasy.dev/util/structs"

	"github.com/allegro/bigcache/v3"
)

type Local struct {
	cache *bigcache.BigCache
}

func NewLocal(ctx context.Context, ttl time.Duration) (Local, error) {
	cacheConfig := bigcache.DefaultConfig(ttl)
	cache, err := bigcache.New(ctx, cacheConfig)
	if err != nil {
		return Local{}, errors.Wrap(err, "unable to create bigcache")
	}

	return Local{
		cache: cache,
	}, nil
}

func (l Local) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := l.cache.Get(key)
	if err == bigcache.ErrEntryNotFound {
		return errors.ErrNotFound
	} else if err != nil {
		return errors.Wrap(err, "unable to retrieve data from local cache")
	}

	err = structs.Unmarshal(data, dest)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal data")
	}

	return nil
}

func (l Local) Put(ctx context.Context, key string, value interface{}) error {
	data, err := structs.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "unable to marshal value")
	}

	err = l.cache.Set(key, data)
	if err != nil {
		return errors.Wrap(err, "unable to set value")
	}

	return nil
}
