package rediscache

import (
	"context"
	"time"

	"goeasy.dev/errors"
	"goeasy.dev/util/structs"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client redis.UniversalClient
	ttl    time.Duration
}

func NewRedisProvider(client redis.UniversalClient, ttl time.Duration) Redis {
	return Redis{
		client: client,
		ttl:    ttl,
	}
}

func (r Redis) Put(ctx context.Context, key string, value interface{}) error {
	data, err := structs.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "unable to marshal value")
	}

	_, err = r.client.Set(ctx, key, data, r.ttl).Result()
	if err != nil {
		return errors.Wrap(err, "unable to store value in redis")
	}

	return nil
}

func (r Redis) Get(ctx context.Context, key string, dest interface{}) error {
	retString, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.ErrNotFound
	} else if err != nil {
		return errors.Wrap(err, "unable to retrieve value from redis")
	}

	err = structs.Unmarshal([]byte(retString), dest)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal value")
	}

	return nil
}
