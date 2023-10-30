package async

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/night-sword/kratos-kit/errors"
)

func TestLoop(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() error
		interval time.Duration
	}{
		{"case.normal", func() error {
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		{"case.error", func() error {
			return errors.NotFound("case.error", time.Now().Format(time.DateTime))
		}, time.Second},

		{"case.panic", func() error {
			panic("case.panic:" + time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		// will 2 second print once
		{"case.slow", func() error {
			time.Sleep(time.Second)
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Loop(tt.fn, tt.interval)
		})
	}

	time.Sleep(time.Second * 10)
}

func TestTick(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() error
		interval time.Duration
	}{
		{"case.normal", func() error {
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		{"case.error", func() error {
			return errors.NotFound("case.error", time.Now().Format(time.DateTime))
		}, time.Second},

		{"case.panic", func() error {
			panic("case.panic:" + time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		// will 1 second print once
		{"case.slow", func() error {
			time.Sleep(time.Second)
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Tick(tt.fn, tt.interval)
		})
	}

	time.Sleep(time.Second * 10)
}

func TestLoopContext(t *testing.T) {
	bg := context.Background()
	timeout := func() context.Context {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		return ctx
	}

	tests := []struct {
		name     string
		ctx      context.Context
		fn       func(ctx context.Context) error
		interval time.Duration
	}{
		{"case.normal", bg, func(ctx context.Context) error {
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		{"case.normal.timeout", timeout(), func(ctx context.Context) error {
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		{"case.error", bg, func(ctx context.Context) error {
			return errors.NotFound("case.error", time.Now().Format(time.DateTime))
		}, time.Second},

		{"case.error.timeout", timeout(), func(ctx context.Context) error {
			return errors.NotFound("case.error", time.Now().Format(time.DateTime))
		}, time.Second},

		{"case.panic", bg, func(ctx context.Context) error {
			panic("case.panic:" + time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		{"case.panic.timeout", timeout(), func(ctx context.Context) error {
			panic("case.panic:" + time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		// will 2 second print once
		{"case.slow", bg, func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		// will 2 second print once
		{"case.slow.timeout", timeout(), func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoopContext(tt.ctx, tt.fn, tt.interval)
		})
	}

	time.Sleep(time.Second * 10)
}

func TestTickContext(t *testing.T) {
	bg := context.Background()
	timeout := func() context.Context {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		return ctx
	}

	tests := []struct {
		name     string
		ctx      context.Context
		fn       func(ctx context.Context) error
		interval time.Duration
	}{
		{"case.normal", bg, func(ctx context.Context) error {
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		{"case.normal.timeout", timeout(), func(ctx context.Context) error {
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		{"case.error", bg, func(ctx context.Context) error {
			return errors.NotFound("case.error", time.Now().Format(time.DateTime))
		}, time.Second},

		{"case.error.timeout", timeout(), func(ctx context.Context) error {
			return errors.NotFound("case.error", time.Now().Format(time.DateTime))
		}, time.Second},

		{"case.panic", bg, func(ctx context.Context) error {
			panic("case.panic:" + time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		{"case.panic.timeout", timeout(), func(ctx context.Context) error {
			panic("case.panic:" + time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		// will 1 second print once
		{"case.slow", bg, func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},

		// will 1 second print once
		{"case.slow.timeout", timeout(), func(ctx context.Context) error {
			time.Sleep(time.Second)
			fmt.Println(time.Now().Format(time.DateTime))
			return nil
		}, time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TickContext(tt.ctx, tt.fn, tt.interval)
		})
	}

	time.Sleep(time.Second * 10)
}
