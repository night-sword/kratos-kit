package log

import (
	"fmt"
	"os"
	"reflect"

	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/night-sword/kratos-kit/cnst"
)

var defaultEncoderCfg = &zapcore.EncoderConfig{
	NameKey:        "logger",
	LevelKey:       LogKeyLevel,
	CallerKey:      LogKeyCaller,
	FunctionKey:    LogKeyFunction,
	StacktraceKey:  LogKeyStack,
	MessageKey:     zapcore.OmitKey,
	SkipLineEnding: false,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var defaultValues = []any{
	LogKeyTimestamp, log.Timestamp("20060102:150405"),
	LogKeyCaller, log.DefaultCaller,
	LogKeyVersion, "",
}

var loggerToNoStack = make(map[log.Logger]log.Logger)

// level : log level
func NewLogger(level string, encoderCfg *zapcore.EncoderConfig, values []any) log.Logger {
	logger := newLogger(level, encoderCfg, zapcore.ErrorLevel, values)
	noStack := newLogger(level, encoderCfg, zapcore.ErrorLevel+1, values)

	loggerToNoStack[logger] = noStack

	return logger
}

func WithNoStack(logger log.Logger) log.Logger {
	noStack, ok := loggerToNoStack[logger]
	if ok {
		return noStack
	}

	return logger
}

func newLogger(outputLv string, encoderCfg *zapcore.EncoderConfig, stackLevel zapcore.Level, values []any) log.Logger {
	z := newZap(outputLv, encoderCfg, stackLevel)
	kz := kzap.NewLogger(z)

	values = mergeValues(defaultValues, values)
	return log.With(kz, values...)
}

func newZap(outputLv string, encoderCfg *zapcore.EncoderConfig, stackLevel zapcore.Level) *zap.Logger {
	cfg := mergeEncoderCfg(defaultEncoderCfg, encoderCfg)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(*cfg),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		formatLevel(outputLv),
	)

	return zap.New(core).WithOptions(
		zap.AddStacktrace(stackLevel),
		zap.AddCallerSkip(3),
	)
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
