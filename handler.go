package solislog

import (
	"io"
	"time"
)

// Handler defines where and how log records are written.
//
// Each handler has its own output writer, minimum level, template,
// time formatting settings, and output mode.
type Handler struct {
	out          io.Writer
	level        Level
	template     []templateSegment
	timeFormat   string
	location     *time.Location
	json         bool
	errorHandler ErrorHandlerFunc
}

// HandlerOptions configures a Handler.
//
// A nil HandlerOptions value passed to NewHandler uses the default template,
// RFC3339 time format, local time zone, and text output mode.
type HandlerOptions struct {
	// Template defines the output fields and their order.
	//
	// In text mode, text and placeholders are rendered as a log line.
	// In JSON mode, placeholders define the JSON fields and their order;
	// plain text parts are ignored.
	Template string

	// TimeFormat defines how the {time} placeholder is formatted.
	//
	// If empty, time.RFC3339 is used.
	TimeFormat string

	// Location defines the time zone used for the {time} placeholder.
	//
	// If nil, time.Local is used.
	Location *time.Location

	// JSON enables JSON output mode.
	JSON bool

	// ErrorHandler handles write errors.
	//
	// If nil, write errors are ignored.
	ErrorHandler ErrorHandlerFunc
}

// ErrorHandlerFunc is called when a log record cannot be written.
//
// The err argument contains the write error returned by the underlying writer.
// The msg argument contains the already rendered log message that failed to write.
//
// ErrorHandlerFunc is optional. If it is nil, write errors are ignored.
type ErrorHandlerFunc func(err error, msg string)

// AddHandler adds a handler to the logger.
//
// The handler is added to the logger's shared core, so it is also used by
// loggers created from this logger with Bind or Contextualize.
func (logger *Logger) AddHandler(handler Handler) {
	logger.core.mutex.Lock()
	defer logger.core.mutex.Unlock()
	logger.core.handlers = append(logger.core.handlers, handler)
}

// NewHandler creates a handler for the given writer and minimum level.
//
// If options is nil, the handler uses the default text template,
// RFC3339 time format, and local time zone.
func NewHandler(out io.Writer, level Level, options *HandlerOptions) Handler {
	if options == nil {
		options = &HandlerOptions{}
	}

	template := options.Template
	if template == "" {
		template = "{time} | {level} | {message}\n"
	}

	timeFormat := options.TimeFormat
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	location := options.Location
	if location == nil {
		location = time.Local
	}

	return Handler{
		out:          out,
		level:        level,
		template:     parseTokens(tokenize(template)),
		timeFormat:   timeFormat,
		location:     location,
		json:         options.JSON,
		errorHandler: options.ErrorHandler,
	}
}
