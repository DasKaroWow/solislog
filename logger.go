package solislog

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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

type callerMetadata struct {
	file     string
	path     string
	line     int
	function string
	caller   string
}

func getCallerMetadata(skip int) callerMetadata {
	pc, path, line, ok := runtime.Caller(skip)
	if !ok {
		return callerMetadata{}
	}

	file := filepath.Base(path)

	function := ""
	if fn := runtime.FuncForPC(pc); fn != nil {
		function = fn.Name()
	}

	return callerMetadata{
		file:     file,
		path:     path,
		line:     line,
		function: function,
		caller:   file + ":" + strconv.Itoa(line),
	}
}

func (logger *Logger) msg(message string, level Level) {
	caller := getCallerMetadata(3)

	currentRecord := &Record{
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

	type errorInfo struct {
		err     error
		msg     string
		handler ErrorHandlerFunc
	}
	type afterHookInfo struct {
		record *Record
		msg    string
		hook   AfterHookFunc
	}
	var errors []errorInfo
	var afterHooks []afterHookInfo

	logger.core.mutex.Lock()
	defer func() {
		for _, info := range afterHooks {
			info.hook(info.record, info.msg)
		}
		for _, info := range errors {
			info.handler(info.err, info.msg)
		}
	}()
	defer logger.core.mutex.Unlock()

	for i := range logger.core.handlers {
		handler := &logger.core.handlers[i]
		if handler.level > level {
			continue
		}

		handlerRecord := currentRecord.clone()
		if handler.beforeHook != nil {
			handler.beforeHook(handlerRecord)
		}

		var rendered string

		if handler.json {
			rendered = renderJSONRecord(handler, handlerRecord)
		} else {
			rendered = renderTemplateRecord(handler, handlerRecord)
		}
		_, err := handler.out.Write([]byte(rendered))

		if err != nil && handler.errorHandler != nil {
			errors = append(errors, errorInfo{err, rendered, handler.errorHandler})
		}
		if handler.afterHook != nil {
			afterHooks = append(afterHooks, afterHookInfo{handlerRecord, rendered, handler.afterHook})
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
