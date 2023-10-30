package log

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultEncoderCfg = &zapcore.EncoderConfig{
	LevelKey:       "LV",
	NameKey:        "logger",
	CallerKey:      "CALLER",
	FunctionKey:    "FUN",
	StacktraceKey:  "STACK",
	SkipLineEnding: false,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var defaultValues = []any{
	"TS", log.Timestamp("20060102:150405"),
	"CALLER", log.DefaultCaller,
	"VER", "",
}

var _zaplogger *zap.Logger

// level : log level
func NewLogger(level string, encoderCfg *zapcore.EncoderConfig, values []any) log.Logger {
	cfg := mergeEncoderCfg(defaultEncoderCfg, encoderCfg)
	values = mergeValues(defaultValues, values)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(*cfg),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		formatLevel(level),
	)
	zlogger := zap.New(core).WithOptions(
		zap.AddStacktrace(zap.ErrorLevel),
	)
	_zaplogger = zlogger
	kratosZapLogger := kzap.NewLogger(zlogger)

	return log.With(kratosZapLogger, values...)
}

// use for some lib bind log lib to zap
func ZapLogger() (logger *zap.Logger, err error) {
	logger = _zaplogger
	if logger == nil {
		err = errors.New("zap logger now init")
	}
	return
}

func formatLevel(level string) zapcore.Level {
	ms := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}

	lv, ok := ms[level]
	if ok {
		return lv
	}
	return zapcore.InfoLevel
}

func mergeValues(vs1, vs2 []any) (values []any) {
	if len(vs2) == 0 {
		return vs1
	}

	if len(vs2)%2 != 0 {
		panic(fmt.Sprint("input values not pairs: ", vs2))
	}

	map1 := make(map[any]any, len(vs1)/2)
	map2 := make(map[any]any, len(vs2)/2)

	for i := 0; i < len(vs1); i += 2 {
		map1[vs1[i]] = vs1[i+1]
	}
	for i := 0; i < len(vs2); i += 2 {
		map2[vs2[i]] = vs2[i+1]
	}

	// merge
	for k, v := range map2 {
		map1[k] = v
	}

	// convert to pairs slice
	for k, v := range map1 {
		values = append(values, k, v)
	}

	return
}

func mergeEncoderCfg(c1, c2 *zapcore.EncoderConfig) *zapcore.EncoderConfig {
	if c2 == nil {
		return c1
	}

	r1 := reflect.ValueOf(c1).Elem()
	r2 := reflect.ValueOf(c2).Elem()

	t1 := r1.Type()

	for i := 0; i < r1.NumField(); i++ {
		v2Field := r2.FieldByName(t1.Field(i).Name)

		if v2Field.IsValid() && !v2Field.IsZero() && r1.Field(i).CanSet() {
			r1.Field(i).Set(v2Field)
		}
	}

	return c1
}
