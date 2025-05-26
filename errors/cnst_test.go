package errors

import (
	"errors"
	"testing"
)

func TestIsUnrecoverable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "unrecoverable error",
			err:  Unrecoverable,
			want: true,
		},
		{
			name: "wrapped unrecoverable error",
			err:  errors.Join(errors.New("some error"), Unrecoverable),
			want: true,
		},
		{
			name: "different error",
			err:  errors.New("some other error"),
			want: false,
		},
		{
			name: "error with similar message but not the same error",
			err:  errors.New(RsnUnrecoverable),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUnrecoverable(tt.err)
			if got != tt.want {
				t.Errorf("IsUnrecoverable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIgnorableErr(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "ignorable error",
			err:  IgnorableErr,
			want: true,
		},
		{
			name: "wrapped ignorable error",
			err:  errors.Join(errors.New("some error"), IgnorableErr),
			want: true,
		},
		{
			name: "different error",
			err:  errors.New("some other error"),
			want: false,
		},
		{
			name: "error with similar message but not the same error",
			err:  errors.New(RsnIgnorable),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIgnorableErr(tt.err)
			if got != tt.want {
				t.Errorf("IsIgnorableErr() = %v, want %v", got, tt.want)
			}
		})
	}
}
