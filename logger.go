package solislog

import (
	"os"
	"sync"
	"time"
)

// Logger writes log messages to one or more handlers.
//
// A Logger keeps default contextual fields in Extra and shares its handler
// configuration with loggers created by Bind or Contextualize.
type Logger struct {
	core  *sharedCore
	extra Extra
}
type sharedCore struct {
	mutex    sync.Mutex
	handlers []Handler
}

// NewLogger creates a logger with default extra fields and handlers.
//
// The provided Extra is copied, so later changes to the original map do not
// affect the logger. Handlers can also be added later with AddHandler.
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
	currentRecord := &record{
		time:    time.Now(),
		level:   level,
		message: message,
		extra:   logger.extra,
	}

	logger.core.mutex.Lock()
	defer logger.core.mutex.Unlock()

	for i := range logger.core.handlers {
		handler := &logger.core.handlers[i]

		if handler.level > level {
			continue
		}

		var rendered string
		if handler.json {
			rendered = renderJSONRecord(handler, currentRecord)
		} else {
			rendered = renderTemplateRecord(handler, currentRecord)
		}
		_, _ = handler.out.Write([]byte(rendered))
	}
}

// Debug logs a message at DebugLevel.
func (logger *Logger) Debug(message string) {
	logger.msg(message, DebugLevel)
}

// Info logs a message at InfoLevel.
func (logger *Logger) Info(message string) {
	logger.msg(message, InfoLevel)
}

// Warning logs a message at WarningLevel.
func (logger *Logger) Warning(message string) {
	logger.msg(message, WarningLevel)
}

// Error logs a message at ErrorLevel.
func (logger *Logger) Error(message string) {
	logger.msg(message, ErrorLevel)
}

// Fatal logs a message at FatalLevel and exits the process with status code 1.
func (logger *Logger) Fatal(message string) {
	logger.msg(message, FatalLevel)
	os.Exit(1)
}
