package log

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	. "github.com/night-sword/kratos-kit/cnst"
	"github.com/night-sword/kratos-kit/errors"
)

// globalLogger is designed as a global logger in current process.
var global = &loggerAppliance{}

// loggerAppliance is the proxy of `Logger` to
// make logger change will affect all sub-logger.
type loggerAppliance struct {
	lock sync.RWMutex
	log.Logger
	messageKey string
}

func init() {
	global.SetLogger(log.DefaultLogger)
	global.SetMessageKey(LogKeyMessage)
}

func (a *loggerAppliance) SetLogger(in log.Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

func (a *loggerAppliance) GetLogger() log.Logger {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.Logger
}

func (a *loggerAppliance) SetMessageKey(key string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.messageKey = key
}

func Unwrap(l log.Logger) log.Logger {
	if u, ok := l.(*loggerAppliance); ok {
		return u.Logger
	}
	return l
}

// SetLogger should be called before any other log call.
// And it is NOT THREAD SAFE.
func SetLogger(logger log.Logger) {
	global.SetLogger(logger)
}

func SetMessageKey(key string) {
	global.SetMessageKey(key)
}

// GetLogger returns global logger appliance as logger in current process.
func GetLogger() log.Logger {
	return global
}

// Log Print log by level and keyvals.
func Log(level log.Level, keyvals ...interface{}) {
	_ = global.Log(level, keyvals...)
}

// Context with context logger.
func Context(ctx context.Context) *log.Helper {
	return log.NewHelper(log.WithContext(ctx, global.Logger))
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	_ = global.Log(log.LevelDebug, global.messageKey, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	_ = global.Log(log.LevelDebug, global.messageKey, fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	_ = global.Log(log.LevelDebug, keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	_ = global.Log(log.LevelInfo, global.messageKey, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	_ = global.Log(log.LevelInfo, global.messageKey, fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func Infow(keyvals ...interface{}) {
	_ = global.Log(log.LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	_ = global.Log(log.LevelWarn, global.messageKey, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	_ = global.Log(log.LevelWarn, global.messageKey, fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func Warnw(keyvals ...interface{}) {
	_ = global.Log(log.LevelWarn, keyvals...)
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	level := log.LevelError
	if err, ok := a[0].(error); ok {
		level = GetLevel(err)
		kvs := ExtractError(err)

		if level == log.LevelError {
			// If the error level is LvError, then instead of using the zap stack, use the error's own stack.
			_ = WithNoStack(Unwrap(global)).Log(level, append(kvs, StackTrace(err)...)...)
			return
		}

		_ = global.Log(level, append(kvs, LogKeyOperation, errors.FromError(err).StackTrace()[0])...)
		return
	}

	_ = global.Log(level, global.messageKey, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	_ = global.Log(log.LevelError, global.messageKey, fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func Errorw(keyvals ...interface{}) {
	_ = global.Log(log.LevelError, keyvals...)
}

// ErrorE logs the error if it is not nil.
func ErrorE(err error) {
	if err != nil {
		Error(err)
	}
}

// E logs the error if it is not nil.
// E is a shorthand version of the ErrorE function.
func E(err error) {
	if err != nil {
		Error(err)
	}
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	_ = global.Log(log.LevelFatal, global.messageKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	_ = global.Log(log.LevelFatal, global.messageKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyvals ...interface{}) {
	_ = global.Log(log.LevelFatal, keyvals...)
	os.Exit(1)
}

func Debugm(msg string, kvs ...any) {
	_ = global.Log(log.LevelDebug, append([]any{global.messageKey, msg}, kvs...)...)
}

func Infom(msg string, kvs ...any) {
	_ = global.Log(log.LevelInfo, append([]any{global.messageKey, msg}, kvs...)...)
}

func Warnm(msg string, kvs ...any) {
	_ = global.Log(log.LevelWarn, append([]any{global.messageKey, msg}, kvs...)...)
}

func Errorm(msg string, kvs ...any) {
	_ = global.Log(log.LevelError, append([]any{global.messageKey, msg}, kvs...)...)
}

func Fatalm(msg string, kvs ...any) {
	_ = global.Log(log.LevelFatal, append([]any{global.messageKey, msg}, kvs...)...)
	os.Exit(1)
}
