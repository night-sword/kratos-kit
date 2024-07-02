package async

import (
	"context"
	"time"

	"github.com/night-sword/kratos-kit/log"
)

func Loop(fn func() error, interval time.Duration) {
	go func() {
		for {
			log.E(Safe(fn))
			time.Sleep(interval)
		}
	}()
}

func LoopContext(ctx context.Context, fn func(ctx context.Context) error, interval time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("context canceled, stopping loop")
				return
			default:
				log.E(SafeContext(ctx, fn))
				time.Sleep(interval)
			}
		}
	}()
}

func Tick(fn func() error, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.E(Safe(fn))
			}
		}
	}()
}

func TickContext(ctx context.Context, fn func(ctx context.Context) error, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info("context canceled, stopping loop")
				return
			case <-ticker.C:
				log.E(SafeContext(ctx, fn))
			}
		}
	}()
}
