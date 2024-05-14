package cnst

const (
	LogKeyAsWarn        = "AsWarn"
	LogKeyUnrecoverable = "Unrecoverable"
	LogKeyLevel         = "LV"
	LogKeyOperation     = "OPT"
	LogKeyTimestamp     = "TS"
	LogKeyCaller        = "CALLER"
	LogKeyVersion       = "VER"
	LogKeyMessage       = "MSG"
	LogKeyLatency       = "LATENCY"
	LogKeyCode          = "CODE"
	LogKeyReason        = "REASON"
	LogKeyMeta          = "META"
	LogKeyCause         = "CAUSE"
	LogKeyStack         = "STACK"
	LogKeyArg           = "ARG"
	LogKeyFunction      = "FUN"
)

var OKValue = "1"

var (
	MetaAsWarn = map[string]string{LogKeyAsWarn: OKValue}
)
