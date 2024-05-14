package log

const (
	KeyAsWarn        = "AsWarn"
	KeyUnrecoverable = "Unrecoverable"
	KeyLevel         = "LV"
	KeyOperation     = "OPT"
	KeyTimestamp     = "TS"
	KeyCaller        = "CALLER"
	KeyVersion       = "VER"
	KeyMessage       = "MSG"
	KeyLatency       = "LATENCY"
	KeyCode          = "CODE"
	KeyReason        = "REASON"
	KeyMeta          = "META"
	KeyCause         = "CAUSE"
	KeyStack         = "STACK"
	KeyArg           = "ARG"
	KeyFunction      = "FUN"
)

var OKValue = "1"

var (
	MetaAsWarn = map[string]string{KeyAsWarn: OKValue}
)
