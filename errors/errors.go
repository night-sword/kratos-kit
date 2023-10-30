package errors

import (
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"

	kerr "github.com/go-kratos/kratos/v2/errors"
	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"

	"github.com/night-sword/kratos-kit/cnst"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

// Error is a status error.
type Error struct {
	kerr.Status
	cause error
	stack *stack
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v cause = %v", e.Code, e.Reason, e.Message, e.Metadata, e.cause)
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.Reason == e.Reason
	}
	return false
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(meta map[string]string) *Error {
	err := Clone(e)
	err.Metadata = meta
	return err
}

func (e *Error) AppendMetadata(k string, v any) *Error {
	if e.Metadata == nil {
		e.Metadata = map[string]string{k: fmt.Sprintf("%v", v)}
	} else {
		e.Metadata[k] = fmt.Sprintf("%v", v)
	}
	return e
}

// Alias AppendMetadata
func (e *Error) AddMetadata(k string, v any) *Error {
	return e.AppendMetadata(k, v)
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(httpstatus.ToGRPCCode(int(e.Code)), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

func (e *Error) Unrecoverable() *Error {
	return e.WithCause(Unrecoverable).AddMetadata(cnst.LogKeyUnrecoverable, cnst.LogOKValue)
}

func (e *Error) AsWarn() *Error {
	return Clone(e).AddMetadata(cnst.LogKeyAsWarn, cnst.LogOKValue)
}

func (e *Error) Degrade() *Error {
	return e.Unrecoverable().AsWarn()
}

func (e *Error) AddAsWarnMeta() *Error {
	return e.AppendMetadata(cnst.LogKeyAsWarn, cnst.LogOKValue)
}

// New returns an error object for the code, message.
func New(code int, reason, message string) *Error {
	return &Error{
		Status: kerr.Status{
			Code:    int32(code),
			Message: message,
			Reason:  reason,
		},
		stack: callers(10),
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, reason, format string, a ...interface{}) *Error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code int, reason, format string, a ...interface{}) error {
	return New(code, reason, fmt.Sprintf(format, a...))
}

// Code returns the http code for an error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200 //nolint:gomnd
	}
	return int(FromError(err).Code)
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// Clone deep clone error to a new error.
func Clone(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &Error{
		Status: kerr.Status{
			Code:     err.Code,
			Reason:   err.Reason,
			Message:  err.Message,
			Metadata: metadata,
		},
		cause: err.cause,
		stack: err.stack,
	}
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if !ok {
		return New(UnknownCode, UnknownReason, err.Error())
	}
	ret := New(
		httpstatus.FromGRPCCode(gs.Code()),
		UnknownReason,
		gs.Message(),
	)
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			ret.Reason = d.Reason
			return ret.WithMetadata(d.Metadata)
		}
	}
	return ret
}

type _rsp struct {
	Code     int32             `json:"code"`
	Reason   string            `json:"reason"`
	Message  string            `json:"message"`
	Metadata map[string]string `json:"metadata"`
}

func FromHttpRsp(body []byte) *Error {
	rsp := new(_rsp)
	if err := json.Unmarshal(body, rsp); err != nil {
		return BadRequest(RsnParams, "decode response to error struct fail")
	}

	err := &Error{
		Status: kerr.Status{
			Code:     rsp.Code,
			Reason:   rsp.Reason,
			Message:  rsp.Message,
			Metadata: rsp.Metadata,
		},
	}

	if u, ok := rsp.Metadata[cnst.LogKeyUnrecoverable]; ok && u == cnst.LogOKValue {
		err.cause = Unrecoverable
	}

	return err
}
