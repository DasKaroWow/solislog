package solislog

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
