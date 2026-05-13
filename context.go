package solislog

import (
	"context"
	"maps"
)

type loggerContextKey struct{}

// Bind returns a new logger with additional contextual fields.
// If extra is empty (nil or zero length), the same logger instance is returned.
//
// The returned logger shares handlers with the original logger. If extra
// contains keys that already exist on the logger, the new values override
// the old ones for the returned logger only.
func (logger *Logger) Bind(extra Extra) *Logger {
	if extra == nil {
		return logger
	}

	merged := maps.Clone(logger.extra)
	maps.Insert(merged, maps.All(extra))

	return &Logger{
		core:  logger.core,
		extra: merged,
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
