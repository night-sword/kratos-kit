package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/night-sword/kratos-kit/errors"
	klog "github.com/night-sword/kratos-kit/log"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// Server is an server logging middleware.
func LogServer(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			startAt := time.Now()
			kvs := make([]any, 0, 24)

			reply, err = handler(ctx, req)

			if info, ok := transport.FromServerContext(ctx); ok {
				kvs = append(kvs, "OPT", info.Operation())
			}
			kvs = append(kvs, "LATENCY", time.Since(startAt).Seconds())

			var level log.Level
			withoutStack := false
			if err != nil {
				level = log.LevelError
				if kerr := errors.FromError(err); kerr != nil {
					withoutStack = true
					if _, asWarn := kerr.GetMetadata()[klog.KeyAsWarn]; asWarn {
						level = log.LevelWarn
					}

					kvs = append(kvs,
						"CODE", kerr.GetCode(),
						"REASON", kerr.GetReason(),
						"MSG", kerr.GetMessage(),
					)
					if len(kerr.GetMetadata()) > 0 {
						kvs = append(kvs, "META", kerr.GetMetadata())
					}
					if kerr.Unwrap() != nil {
						kvs = append(kvs, "CAUSE", kerr.Unwrap())
					}
					if level >= log.LevelError {
						kvs = append(kvs, "STACK", kerr.StackTrace())
					}
				} else {
					kvs = append(kvs, "CODE", errors.UnknownCode, "MSG", err.Error())
				}
			} else {
				level = log.LevelInfo
				kvs = append(kvs, "CODE", 0)
			}

			kvs = append(kvs, "ARGS", extractArgs(req))

			if withoutStack {
				logger = klog.WithNoStack(klog.Unwrap(logger))
			}
			_ = log.WithContext(ctx, logger).Log(level, kvs...)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req any) any {
	if j, err := json.Marshal(req); err == nil {
		var v map[string]any
		if err = json.Unmarshal(j, &v); err == nil {
			return v
		}
	}

	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}

	return fmt.Sprintf("%+v", req)
}
