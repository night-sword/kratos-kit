package log

import (
	"reflect"
	"testing"

	"go.uber.org/zap/zapcore"
)

func Test_mergeValues(t *testing.T) {
	type args struct {
		vs1 []any
		vs2 []any
	}
	tests := []struct {
		name       string
		args       args
		wantValues []any
	}{
		{
			name: "vs2 empty",
			args: args{
				vs1: []any{"a", 1, "b", 2},
				vs2: []any{},
			},
			wantValues: []any{"a", 1, "b", 2},
		},
		{
			name: "vs2 same as vs1",
			args: args{
				vs1: []any{"a", 1, "b", 2},
				vs2: []any{"a", 1, "b", 2},
			},
			wantValues: []any{"a", 1, "b", 2},
		},
		{
			name: "vs2 cover vs1",
			args: args{
				vs1: []any{"a", 1, "b", 2, "c", 3},
				vs2: []any{"a", 11, "b", 22},
			},
			wantValues: []any{"a", 11, "b", 22, "c", 3},
		},
		{
			name: "vs2 new key vs1",
			args: args{
				vs1: []any{"a", 1, "b", 2, "c", 3},
				vs2: []any{"d", 4, "b", 22},
			},
			wantValues: []any{"a", 1, "b", 22, "c", 3, "d", 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotValues := mergeValues(tt.args.vs1, tt.args.vs2); !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("mergeValues() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}

func Test_mergeZapCfg(t *testing.T) {
	type args struct {
		c1 *zapcore.EncoderConfig
		c2 *zapcore.EncoderConfig
	}
	tests := []struct {
		name string
		args args
		want *zapcore.EncoderConfig
	}{
		{
			name: "v2 replace v1",
			args: args{
				c1: &zapcore.EncoderConfig{
					MessageKey: "KEY_1",
					LevelKey:   "LV_1",
				},
				c2: &zapcore.EncoderConfig{
					MessageKey: "KEY_2",
					LevelKey:   "LV_2",
				},
			},
			want: &zapcore.EncoderConfig{
				MessageKey: "KEY_2",
				LevelKey:   "LV_2",
			},
		},
		{
			name: "v2 merge v1",
			args: args{
				c1: &zapcore.EncoderConfig{
					MessageKey: "KEY_1",
					LevelKey:   "LV_1",
				},
				c2: &zapcore.EncoderConfig{
					MessageKey: "KEY_2",
					LevelKey:   "LV_1",
					NameKey:    "N_1",
				},
			},
			want: &zapcore.EncoderConfig{
				MessageKey: "KEY_2",
				LevelKey:   "LV_1",
				NameKey:    "N_1",
			},
		},
		{
			name: "v2 merge v1",
			args: args{
				c1: &zapcore.EncoderConfig{
					MessageKey: "KEY_1",
				},
				c2: &zapcore.EncoderConfig{
					MessageKey: "KEY_2",
					LevelKey:   "LV_2",
				},
			},
			want: &zapcore.EncoderConfig{
				MessageKey: "KEY_2",
				LevelKey:   "LV_2",
			},
		},
		{
			name: "v2 zero field not cover v1",
			args: args{
				c1: &zapcore.EncoderConfig{
					MessageKey: "KEY_1",
					LevelKey:   "LV_1",
				},
				c2: &zapcore.EncoderConfig{
					MessageKey: "KEY_2",
				},
			},
			want: &zapcore.EncoderConfig{
				MessageKey: "KEY_2",
				LevelKey:   "LV_1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeEncoderCfg(tt.args.c1, tt.args.c2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeEncoderCfg() = %v, want %v", got, tt.want)
			}
		})
	}
}
