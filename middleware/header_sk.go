package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/night-sword/kratos-kit/errors"
)

func HeaderSK(ak, sk string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			info, ok := transport.FromServerContext(ctx)
			if !ok {
				err = errors.Forbidden(errors.RsnForbidden, "transport parse fail").Degrade()
				return
			}

			if !_isHealthCheckOperation(info) {
				if info.RequestHeader().Get(ak) != sk {
					err = errors.Forbidden(errors.RsnForbidden, "access forbidden").Degrade().AddMetadata("operation", info.Operation())
					return
				}
			}

			return handler(ctx, req)
		}
	}
}

func _isHealthCheckOperation(t transport.Transporter) bool {
	ss := strings.Split(t.Operation(), "/")
	return "HealthCheck" == ss[len(ss)-1]
}
