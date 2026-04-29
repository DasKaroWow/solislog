package solislog

// Level represents the severity of a log message.
type Level int

const (
	// DebugLevel is used for detailed diagnostic messages.
	DebugLevel Level = iota

	// InfoLevel is used for general informational messages.
	InfoLevel

	// WarningLevel is used for messages about unexpected but non-fatal situations.
	WarningLevel

	// ErrorLevel is used for errors that should be visible to the caller.
	ErrorLevel

	// FatalLevel is used for fatal errors.
	//
	// Logger.Fatal logs the message and then exits the process with status code 1.
	FatalLevel
)

// String returns the uppercase text representation of the level.
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
