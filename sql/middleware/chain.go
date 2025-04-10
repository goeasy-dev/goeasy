package middleware

import (
	"context"
)

type SQLRequest struct {
	query   string
	args    []interface{}
	context context.Context
}

type Chain []DBMiddleware

type SQLRequestFunc func(SQLRequest) error

type DBMiddleware func(next SQLRequestFunc) SQLRequestFunc

func (c Chain) Exec(ctx context.Context, query string, args []interface{}, handler func(ctx context.Context) error) error {
	r := SQLRequest{
		context: ctx,
		query:   query,
		args:    args,
	}

	f := func(request SQLRequest) error {
		return handler(request.context)
	}

	for i := len(c) - 1; i >= 0; i-- {
		f = c[i](f)
	}

	return f(r)
}
