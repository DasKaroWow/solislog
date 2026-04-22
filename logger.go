package solislog

import (
	"context"
	"io"
	"maps"
	"time"
)

type Logger struct {
	out      io.Writer
	level    Level
	template []templatePart
	extra    map[string]string
}

func Add(out io.Writer, level Level, template string, defaultExtra map[string]string) *Logger { // filter by lambda func by record
	logger := new(Logger)
	logger.out = out
	logger.level = level
	logger.template = parseTemplate(template)
	if defaultExtra != nil {
		logger.extra = maps.Clone(defaultExtra)
	} else {
		logger.extra = map[string]string{}
	}

	return logger
}

type loggerContextKey struct{}

func cloneMapWithDefault(src map[string]string, defaultExtra map[string]string) map[string]string {
	dst := maps.Clone(defaultExtra)
	maps.Insert(dst, maps.All(src))
	return dst
}

func (logger *Logger) Contextualize(ctx context.Context, extra map[string]string) context.Context {
	contextualLogger := &Logger{
		out:      logger.out,
		level:    logger.level,
		template: logger.template,
		extra:    cloneMapWithDefault(extra, logger.extra),
	}
	return context.WithValue(ctx, loggerContextKey{}, contextualLogger)
}

func FromContext(ctx context.Context) (*Logger, bool) {
	logger, ok := ctx.Value(loggerContextKey{}).(*Logger)
	return logger, ok
}

func (logger *Logger) msg(message string, level Level) error {
	currentRecord := new(record)
	currentRecord.time = time.Now()
	currentRecord.level = level
	currentRecord.message = message
	currentRecord.extra = logger.extra
	if logger.level > currentRecord.level {
		return nil
	}
	rendered := renderRecord(logger.template, currentRecord)

	_, err := logger.out.Write([]byte(rendered))
	return err
}

func (logger *Logger) Debug(message string) error {
	return logger.msg(message, DebugLevel)
}

func (logger *Logger) Info(message string) error {
	return logger.msg(message, InfoLevel)
}

func (logger *Logger) Warning(message string) error {
	return logger.msg(message, WarningLevel)
}

func (logger *Logger) Error(message string) error {
	return logger.msg(message, ErrorLevel)
}
