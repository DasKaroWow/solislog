package solislog

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

const ansiReset = "\x1b[0m"

var ansiColors = map[string]string{
	"black":   "\x1b[30m",
	"red":     "\x1b[31m",
	"green":   "\x1b[32m",
	"yellow":  "\x1b[33m",
	"blue":    "\x1b[34m",
	"magenta": "\x1b[35m",
	"cyan":    "\x1b[36m",
	"white":   "\x1b[37m",
	"gray":    "\x1b[90m",
}

type Record struct {
	// Time is the moment when the log record was created.
	Time time.Time

	// Level is the severity of the log record.
	Level Level

	// Extra contains contextual key-value fields attached to the logger.
	Extra Extra

	// Message is the log message text passed to the logger.
	Message string

	// File is the base name of the source file where the log call was made.
	File string

	// Path is the full source file path where the log call was made.
	Path string

	// Line is the source line number where the log call was made.
	Line int

	// Function is the full function name where the log call was made.
	Function string

	// Caller is the compact source location in the form "file:line".
	Caller string
}

func (rec *Record) clone() *Record {
	cloned := *rec
	cloned.Extra = cloneExtra(rec.Extra)
	return &cloned
}

func renderField(handler *Handler, rec *Record, segment *templateSegment) string {
	switch segment.value {
	case "time":
		return rec.Time.In(handler.location).Format(handler.timeFormat)
	case "level":
		return rec.Level.String()
	case "message":
		return rec.Message
	case "extra":
		data, err := json.Marshal(rec.Extra)
		if err != nil {
			return "{}"
		}
		return string(data)
	case "file":
		return rec.File
	case "path":
		return rec.Path
	case "line":
		if rec.Line == 0 {
			return ""
		}
		return strconv.Itoa(rec.Line)
	case "function":
		return rec.Function
	case "caller":
		return rec.Caller
	}
	return ""
}

func renderSegment(handler *Handler, rec *Record, segment *templateSegment) string {
	switch segment.mode {
	case fieldMode:
		return renderField(handler, rec, segment)
	case extraMode:
		return rec.Extra[segment.value]
	}

	return segment.value
}

func renderTemplateRecord(handler *Handler, rec *Record) string {
	var renderedRecord strings.Builder

	renderColor := func(colorName string, level Level) string {
		switch colorName {
		case "level":
			return level.ansiCode()
		case "":
			return ansiReset
		}
		return ansiColors[colorName]

	}

	previousColor := ""

	for i := range handler.template {
		segment := &handler.template[i]

		if segment.color != previousColor {
			renderedRecord.WriteString(renderColor(segment.color, rec.Level))
			previousColor = segment.color
		}
		renderedRecord.WriteString(renderSegment(handler, rec, segment))
	}

	if previousColor != "" {
		renderedRecord.WriteString(ansiReset)
	}
	return renderedRecord.String()
}

func renderJSONRecord(handler *Handler, rec *Record) string {
	var rendered strings.Builder
	rendered.WriteByte('{')

	written := 0

	for i := range handler.template {
		segment := &handler.template[i]
		if segment.mode == textMode {
			continue
		}

		if written > 0 {
			rendered.WriteByte(',')
		}

		renderedSegment := renderSegment(handler, rec, segment)

		name, _ := json.Marshal(segment.value)
		rendered.Write(name)
		rendered.WriteByte(':')

		value, _ := json.Marshal(renderedSegment)
		if segment.mode == fieldMode && segment.value == "extra" {
			value = []byte(renderedSegment)
		}
		rendered.Write(value)
		written++
	}

	rendered.WriteString("}\n")
	return rendered.String()
}
