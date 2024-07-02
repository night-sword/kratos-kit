package errors

import (
	"errors"
	"strings"

	"github.com/avast/retry-go/v4"
)

// ---- Reason ---- //
const (
	RsnNotFound         string = "NOT_FOUND"
	RsnParams                  = "PARAMS_ERROR"
	RsnUnrecoverable           = "UNRECOVERABLE"
	RsnInternal                = "INTERNAL"
	RsnForbidden               = "FORBIDDEN"
	RsnConflict                = "CONFLICT"
	RsnRequestRateLimit        = "REQUEST_RATE_LIMIT"
	RsnAccessRepoFail          = "ACCESS_REPO_FAIL"
	RsnTimeout                 = "TIMEOUT"
	RsnPanic                   = "PANIC"
)

// ---- Unrecoverable ---- //
var Unrecoverable = retry.Unrecoverable(errors.New(RsnUnrecoverable))

func IsUnrecoverable(err error) bool {
	return errors.Is(err, Unrecoverable)
}

// ---- SQL Duplicate ---- //

func IsDuplicateErr(err error) (is bool) {
	if err == nil {
		return
	}

	msg := strings.ToLower(err.Error())

	idx1 := strings.Index(msg, "error 1062")
	idx2 := strings.Index(msg, "duplicate entry")

	is = idx1 >= 0 && idx2 >= 0

	return
}
