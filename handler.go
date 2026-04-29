package solislog

import (
	"io"
	"time"
)

type Handler struct {
	out        io.Writer
	level      Level
	template   []templatePart
	timeFormat string
	location   *time.Location
	json       bool
}

type HandlerOptions struct {
	Template   string
	TimeFormat string
	Location   *time.Location
	JSON       bool
}

func (logger *Logger) AddHandler(handler Handler) {
	logger.core.mutex.Lock()
	defer logger.core.mutex.Unlock()
	logger.core.handlers = append(logger.core.handlers, handler)
}

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
		out:        out,
		level:      level,
		template:   parseTemplate(template),
		timeFormat: timeFormat,
		location:   location,
		json:       options.JSON,
	}
}
