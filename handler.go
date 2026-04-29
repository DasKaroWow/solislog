package solislog

import (
	"io"
)

type Handler struct {
	out      io.Writer
	level    Level
	template []templatePart
}

func (logger *Logger) AddHandler(handler Handler) {
	logger.core.mutex.Lock()
	defer logger.core.mutex.Unlock()
	logger.core.handlers = append(logger.core.handlers, handler)
}

func NewHandler(out io.Writer, level Level, template string) Handler {
	return Handler{
		out:      out,
		level:    level,
		template: parseTemplate(template),
	}
}
