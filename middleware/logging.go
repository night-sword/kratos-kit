package middleware

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/night-sword/kratos-kit/cnst"
	"github.com/night-sword/kratos-kit/errors"
	. "github.com/night-sword/kratos-kit/log"
)

// Server is an server logging middleware.
func LogServer(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			startAt := time.Now()
			kvs := make([]any, 0, 24)

			reply, err = handler(ctx, req)

			if info, ok := transport.FromServerContext(ctx); ok {
				kvs = append(kvs, cnst.LogKeyOperation, info.Operation())
			}
			kvs = append(kvs, cnst.LogKeyLatency, time.Since(startAt).Seconds())

			var level log.Level
			withoutStack := false
			if err != nil {
				level = log.LevelError

				if kerr := errors.FromError(err); kerr != nil {
					withoutStack = true
					level = GetLevel(kerr)
					kvs = append(kvs, ExtractError(kerr)...)

					if level >= log.LevelError {
						kvs = append(kvs, cnst.LogKeyStack, kerr.StackTrace())
					}
				} else {
					kvs = append(kvs, cnst.LogKeyCode, errors.UnknownCode, cnst.LogKeyMessage, err.Error())
				}
			} else {
				level = log.LevelInfo
				kvs = append(kvs, cnst.LogKeyCode, 0)
			}

			kvs = append(kvs, cnst.LogKeyArg, ExtractArgs(req))

			if withoutStack {
				logger = WithNoStack(Unwrap(logger))
			}
			_ = log.WithContext(ctx, logger).Log(level, kvs...)
			return
		}
	}
}
