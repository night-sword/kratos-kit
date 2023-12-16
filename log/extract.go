package log

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/night-sword/kratos-kit/errors"
)

func GetLevel(err error) (lv log.Level) {
	lv = log.LevelError
	kerr := errors.FromError(err)
	if _, asWarn := kerr.GetMetadata()[KeyAsWarn]; asWarn {
		lv = log.LevelWarn
	}
	return
}

func ExtractError(err error) (kvs []any) {
	kerr := errors.FromError(err)

	kvs = append(kvs,
		KeyCode, kerr.GetCode(),
		KeyReason, kerr.GetReason(),
		KeyMessage, kerr.GetMessage(),
	)

	if len(kerr.GetMetadata()) > 0 {
		kvs = append(kvs, KeyMeta, kerr.GetMetadata())
	}

	if kerr.Unwrap() != nil {
		kvs = append(kvs, KeyCause, kerr.Unwrap())
	}

	return
}
