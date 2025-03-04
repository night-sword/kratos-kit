package util

import (
	"context"

	"github.com/avast/retry-go/v4"
)

func Do(ctx context.Context, fn func(ctx context.Context) error, opts ...retry.Option) (err error) {
	if len(opts) == 0 {
		return fn(ctx)
	}

	opts = append(opts, retry.Context(ctx))
	return retry.Do(func() error {
		return fn(ctx)
	}, opts...)
}

func DoResult[T any](ctx context.Context, fn func(ctx context.Context) (T, error), opts ...retry.Option) (rsp T, err error) {
	if len(opts) == 0 {
		return fn(ctx)
	}

	opts = append(opts, retry.Context(ctx))
	return retry.DoWithData[T](func() (T, error) {
		return fn(ctx)
	}, opts...)
}
