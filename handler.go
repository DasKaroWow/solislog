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
	out      io.Writer
	level    Level
	template []templateSegment
	options  HandlerOptions
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

	// WithCaller enables caller fields in template.
	WithCaller bool

	// ErrorHandler handles write errors.
	//
	// If nil, write errors are ignored.
	ErrorHandler ErrorHandlerFunc

	// BeforeHook is called before rendering a log record.
	//
	// NOTE: The record is passed by reference. Modifying it affects ONLY this handler's output.
	// If this handler has no BeforeHook, the record is shared (read-only) for performance.
	BeforeHook BeforeHookFunc

	// AfterHook is called after rendering a log record.
	//
	// It receives the record and the rendered output.
	AfterHook AfterHookFunc

	// ReadOnly skips record cloning. Mutating the record in this mode breaks handler isolation.
	ReadOnly bool
}

// ErrorHandlerFunc is called when a log record cannot be written.
//
// The err argument contains the write error returned by the underlying writer.
// The msg argument contains the already rendered log message that failed to write.
//
// ErrorHandlerFunc is optional. If it is nil, write errors are ignored.
type ErrorHandlerFunc func(record *Record, msg []byte, err error)

// BeforeHookFunc is called before a record is rendered.
//
// It can modify the record, for example by changing Message or adding Extra fields.
type BeforeHookFunc func(record *Record)

// AfterHookFunc is called after a record is rendered.
//
// The rendered argument contains the final rendered output written by the handler.
type AfterHookFunc func(record *Record, msg []byte, successful bool)

// AddHandler adds a handler to the logger.
//
// The handler is added to the logger's shared core, so it is also used by
// loggers created from this logger with Bind or Contextualize.
func (logger *Logger) AddHandler(handler Handler) {
	for {
		oldHandlersPtr := logger.core.handlers.Load()

		newHandlers := make([]Handler, len(*oldHandlersPtr)+1)
		copy(newHandlers, *oldHandlersPtr)
		newHandlers[len(*oldHandlersPtr)] = handler

		newMin := FatalLevel
		newWithCaller := false
		for _, h := range newHandlers {
			newMin = min(h.level, newMin)
			newWithCaller = newWithCaller || h.options.WithCaller
		}

		if logger.core.handlers.CompareAndSwap(oldHandlersPtr, &newHandlers) {
			logger.core.minLevel.Store(int32(newMin))
			logger.core.withCaller.Store(newWithCaller)
			return
		}
	}
}

// NewHandler creates a handler for the given writer and minimum level.
//
// If options is nil, the handler uses the default text template,
// RFC3339 time format, and local time zone.
func NewHandler(out io.Writer, level Level, options *HandlerOptions) Handler {
	if options == nil {
		options = &HandlerOptions{}
	}

	if options.Template == "" {
		options.Template = "{time} | {level} | {message}\n"
	}

	if options.TimeFormat == "" {
		options.TimeFormat = time.RFC3339
	}

	if options.Location == nil {
		options.Location = time.Local
	}

	return Handler{
		out:      out,
		level:    level,
		template: buildSegments(scanTemplate(options.Template)),
		options:  *options,
	}
}
