package solislog

import (
	"os"
	"time"
)

type Logger struct {
	core  *sharedCore
	extra Extra
}
type sharedCore struct {
	handlers []Handler
}

// NewLogger creates a base logger.
//
// A logger stores extra fields and a shared handler core.
// Each handler defines its own writer, level and template.
func NewLogger(defaultExtra Extra, handlers ...Handler) *Logger {
	logger := new(Logger)
	logger.core = new(sharedCore)
	logger.core.handlers = []Handler{}
	logger.extra = cloneExtra(defaultExtra)

	for _, handler := range handlers {
		logger.AddHandler(handler)
	}

	return logger
}

func (logger *Logger) msg(message string, level Level) {
	currentRecord := record{
		time:    time.Now(),
		level:   level,
		message: message,
		extra:   logger.extra,
	}

	for _, handler := range logger.core.handlers {
		if handler.level > level {
			continue
		}

		rendered := renderRecord(handler.template, &currentRecord)
		_, _ = handler.out.Write([]byte(rendered))
	}
}

func (logger *Logger) Debug(message string) {
	logger.msg(message, DebugLevel)
}

func (logger *Logger) Info(message string) {
	logger.msg(message, InfoLevel)
}

func (logger *Logger) Warning(message string) {
	logger.msg(message, WarningLevel)
}

func (logger *Logger) Error(message string) {
	logger.msg(message, ErrorLevel)
}

func (logger *Logger) Fatal(message string) {
	logger.msg(message, FatalLevel)
	os.Exit(1)
}
