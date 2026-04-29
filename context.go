package solislog

import (
	"context"
)

type loggerContextKey struct{}

func (logger *Logger) Bind(extra Extra) *Logger {
	return &Logger{
		core:  logger.core,
		extra: mergeExtra(logger.extra, extra),
	}
}

func (logger *Logger) Contextualize(ctx context.Context, extra Extra) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger.Bind(extra))
}

func FromContext(ctx context.Context) (*Logger, bool) {
	logger, ok := ctx.Value(loggerContextKey{}).(*Logger)
	return logger, ok
}
