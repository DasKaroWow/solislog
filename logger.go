package solislog

import (
	"maps"
	"os"
	"sync/atomic"
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
	handlers   atomic.Pointer[[]Handler]
	minLevel   atomic.Int32
	withCaller atomic.Bool
}

// NewLogger creates a logger with default extra fields and handlers.
//
// The provided Extra is copied, so later changes to the original map do not
// affect the logger. Handlers can also be added later with AddHandler.
func NewLogger(defaultExtra Extra, handlers ...Handler) *Logger {
	extra := maps.Clone(defaultExtra)
	if extra == nil {
		extra = Extra{}
	}
	logger := &Logger{
		core:  new(sharedCore),
		extra: extra,
	}
	logger.core.handlers.Store(&[]Handler{})
	logger.core.minLevel.Store(int32(FatalLevel))

	for _, handler := range handlers {
		logger.AddHandler(handler)
	}

	return logger
}

func (logger *Logger) msg(message string, level Level) {
	if Level(logger.core.minLevel.Load()) > level {
		return
	}
	caller := callerMetadata{}
	if logger.core.withCaller.Load() {
		caller = getCallerMetadata(3)
	}

	record := Record{
		Time:     time.Now(),
		Level:    level,
		Message:  message,
		Extra:    logger.extra,
		File:     caller.file,
		Path:     caller.path,
		Line:     caller.line,
		Function: caller.function,
		Caller:   caller.caller,
	}

	handlers := *logger.core.handlers.Load()
	for i := range handlers {
		handler := &handlers[i]
		if handler.level > level {
			continue
		}
		currentRecord := &record
		if !handler.options.ReadOnly && (handler.options.BeforeHook != nil || handler.options.AfterHook != nil || handler.options.ErrorHandler != nil) {
			cloned := record
			cloned.Extra = maps.Clone(record.Extra)
			currentRecord = &cloned
		}
		if handler.options.BeforeHook != nil {
			handler.options.BeforeHook(currentRecord)
		}

		var rendered []byte
		if handler.options.JSON {
			rendered = renderJSONRecord(handler, currentRecord)
		} else {
			rendered = renderTemplateRecord(handler, currentRecord)
		}
		_, err := handler.out.Write(rendered)

		if err != nil && handler.options.ErrorHandler != nil {
			handler.options.ErrorHandler(currentRecord, rendered, err)
		}
		if handler.options.AfterHook != nil {
			handler.options.AfterHook(currentRecord, rendered, err == nil)
		}
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
