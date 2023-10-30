package log

import (
	"encoding/json"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

	. "github.com/night-sword/kratos-kit/cnst"
	"github.com/night-sword/kratos-kit/errors"
)

func GetLevel(err error) (lv log.Level) {
	lv = log.LevelError
	kerr := errors.FromError(err)
	if _, asWarn := kerr.GetMetadata()[LogKeyAsWarn]; asWarn {
		lv = log.LevelWarn
	}
	return
}

func ExtractError(err error) (kvs []any) {
	kerr := errors.FromError(err)

	kvs = append(kvs,
		LogKeyCode, kerr.GetCode(),
		LogKeyReason, kerr.GetReason(),
		LogKeyMessage, kerr.GetMessage(),
	)

	if len(kerr.GetMetadata()) > 0 {
		kvs = append(kvs, LogKeyMeta, kerr.GetMetadata())
	}

	if kerr.Unwrap() != nil {
		kvs = append(kvs, LogKeyCause, kerr.Unwrap())
	}

	return
}

func StackTrace(err error) (kvs []any) {
	kerr := errors.FromError(err)

	return []any{LogKeyStack, kerr.StackTrace()}
}

// extractArgs returns the string of the req
func ExtractArgs(req any) any {
	if j, err := json.Marshal(req); err == nil {
		var v map[string]any
		if err = json.Unmarshal(j, &v); err == nil {
			return v
		}
	}

	if _redacter, ok := req.(Redacter); ok {
		return _redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}

	return fmt.Sprintf("%+v", req)
}

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}
