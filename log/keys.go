package log

const (
	KeyAsWarn    string = "AsWarn"
	KeyLevel            = "LV"
	KeyOperation        = "OPT"
	KeyTimestamp        = "TS"
	KeyCaller           = "CALLER"
	KeyVersion          = "VER"
	KeyMessage          = "MSG"
	KeyLatency          = "LATENCY"
	KeyCode             = "CODE"
	KeyReason           = "REASON"
	KeyMeta             = "META"
	KeyCause            = "CAUSE"
	KeyStack            = "STACK"
	KeyArg              = "ARG"
	KeyFunction         = "FUN"
)

var (
	MetaAsWarn = map[string]string{KeyAsWarn: "1"}
)
