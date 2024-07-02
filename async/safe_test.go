package async

import (
	"context"
	"testing"
	"time"

	"github.com/night-sword/kratos-kit/errors"
)

func TestSafe(t *testing.T) {
	tests := []struct {
		name    string
		fn      func() error
		wantErr bool
	}{
		{"case.1", func() error { return nil }, false},
		{"case.2", func() error {
			panic(111)
			return nil
		}, true},
		{"case.3", func() error { return errors.NotFound("1", "2") }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Safe(tt.fn); (err != nil) != tt.wantErr {
				t.Errorf("Safe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSafeContext(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	cancel()

	tests := []struct {
		name    string
		ctx     context.Context
		fn      func(ctx context.Context) error
		wantErr bool
	}{
		{"case.1", context.Background(), func(ctx context.Context) error { return nil }, false},
		{"case.2", context.Background(), func(ctx context.Context) error {
			panic(111)
			return nil
		}, true},
		{"case.3", ctxTimeout, func(ctx context.Context) error { return errors.NotFound("1", "2") }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SafeContext(tt.ctx, tt.fn); (err != nil) != tt.wantErr {
				t.Errorf("SafeContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
