package service

import (
	"context"

	"goeasy.dev/cache"
)

type Service struct{}

func (*Service) Get(ctx context.Context, key string) (dest string, err error) {
	err = cache.Get(ctx, key, &dest)
	return
}

func (*Service) Put(ctx context.Context, key, value string) error {
	return cache.Put(ctx, key, value)
}
