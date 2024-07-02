package async

import (
	"context"
	"fmt"

	"github.com/night-sword/kratos-kit/errors"
)

func SafeContext(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = errors.InternalServer(errors.RsnPanic, fmt.Sprintf("%v", p))
		}
	}()

	return fn(ctx)
}

func Safe(fn func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = errors.InternalServer(errors.RsnPanic, fmt.Sprintf("%v", p))
		}
	}()

	return fn()
}
