package log

import (
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/log"

	. "github.com/night-sword/kratos-kit/cnst"
	"github.com/night-sword/kratos-kit/errors"
)

func TestGetLevel(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		args   args
		wantLv log.Level
	}{
		{"case.1", args{errors.BadRequest("1", "")}, log.LevelError},
		{"case.2", args{errors.BadRequest("1", "").WithMetadata(MetaAsWarn)}, log.LevelWarn},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLv := GetLevel(tt.args.err); gotLv != tt.wantLv {
				t.Errorf("GetLevel() = %v, want %v", gotLv, tt.wantLv)
			}
		})
	}
}

func TestExtractError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantKvs []any
	}{
		{"case.1", args{errors.BadRequest("1", "2").WithMetadata(MetaAsWarn)}, []any{LogKeyCode, errors.BadRequest("1", "1").GetCode(), LogKeyReason, "1", LogKeyMessage, "2", LogKeyMeta, MetaAsWarn}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotKvs := ExtractError(tt.args.err); !reflect.DeepEqual(gotKvs, tt.wantKvs) {
				t.Errorf("ExtractError() = %v, want %v", gotKvs, tt.wantKvs)
			}
		})
	}
}
