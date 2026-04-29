package solislog

import (
	"context"
)

type loggerContextKey struct{}

// Bind returns a new logger with additional contextual fields.
//
// The returned logger shares handlers with the original logger. If extra
// contains keys that already exist on the logger, the new values override
// the old ones for the returned logger only.
func (logger *Logger) Bind(extra Extra) *Logger {
	return &Logger{
		core:  logger.core,
		extra: mergeExtra(logger.extra, extra),
	}
}

// Contextualize returns a context containing a bound logger.
//
// The stored logger shares handlers with the original logger and includes
// the provided extra fields.
func (logger *Logger) Contextualize(ctx context.Context, extra Extra) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger.Bind(extra))
}

// FromContext returns a logger stored in ctx by Contextualize.
//
// The second return value reports whether a logger was found.
func FromContext(ctx context.Context) (*Logger, bool) {
	logger, ok := ctx.Value(loggerContextKey{}).(*Logger)
	return logger, ok
}
