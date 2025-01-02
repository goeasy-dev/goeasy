package cache

import (
	"context"

	"goeasy.dev/errors"
)

type noopProvider struct{}

func (n *noopProvider) Put(ctx context.Context, key string, value interface{}) error {
	return nil
}

func (n *noopProvider) Get(ctx context.Context, key string, dest interface{}) error {
	return errors.ErrNotFound
}
