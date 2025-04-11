package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/night-sword/kratos-kit/errors"
)

func FormatError() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if reply, err = handler(ctx, req); err != nil {
				err = errors.FromError(err)
			}
			return

		}
	}
}
