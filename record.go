package solislog

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
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

// Record contains all data used to render a single log entry.
//
// A Record is passed to hooks before and after rendering. Before hooks may
// modify the record before it is rendered by a handler.
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

	buf := make([]byte, 0, len(file)+12)
	buf = append(buf, file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(line), 10)

	return callerMetadata{
		file:     file,
		path:     path,
		line:     line,
		function: function,
		caller:   string(buf),
	}
}

type fieldResolver func(handler *Handler, record *Record) string

var fieldResolvers = map[string]fieldResolver{
	"time": func(handler *Handler, record *Record) string {
		return record.Time.In(handler.options.Location).Format(handler.options.TimeFormat)
	},
	"level": func(_ *Handler, record *Record) string {
		return record.Level.String()
	},
	"message": func(_ *Handler, record *Record) string {
		return record.Message
	},
	"extra": func(_ *Handler, record *Record) string {
		data, err := json.Marshal(record.Extra)
		if err != nil {
			return "{}"
		}
		return string(data)
	},
	"file": func(_ *Handler, record *Record) string {
		return record.File
	},
	"path": func(_ *Handler, record *Record) string {
		return record.Path
	},
	"line": func(_ *Handler, record *Record) string {
		if record.Line == 0 {
			return ""
		}
		return strconv.Itoa(record.Line)
	},
	"function": func(_ *Handler, record *Record) string {
		return record.Function
	},
	"caller": func(_ *Handler, record *Record) string {
		return record.Caller
	},
}

func renderSegment(handler *Handler, rec *Record, segment *templateSegment) string {
	switch segment.mode {
	case fieldMode:
		if resolver, ok := fieldResolvers[segment.value]; ok {
			return resolver(handler, rec)
		}
		return ""
	case extraMode:
		return rec.Extra[segment.value]
	}
	return segment.value
}

var bufferPool = sync.Pool{
	New: func() any {
		buf := new(bytes.Buffer)
		buf.Grow(512)
		return buf
	},
}

func renderTemplateRecord(handler *Handler, rec *Record) []byte {
	renderedRecord := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(renderedRecord)
	renderedRecord.Reset()

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

	result := make([]byte, renderedRecord.Len())
	copy(result, renderedRecord.Bytes())
	return result
}

func renderJSONRecord(handler *Handler, rec *Record) []byte {
	renderedRecord := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(renderedRecord)
	renderedRecord.Reset()
	renderedRecord.WriteByte('{')

	written := 0

	for i := range handler.template {
		segment := &handler.template[i]
		if segment.mode == textMode {
			continue
		}

		if written > 0 {
			renderedRecord.WriteByte(',')
		}

		renderedSegment := renderSegment(handler, rec, segment)

		name, _ := json.Marshal(segment.value)
		renderedRecord.Write(name)
		renderedRecord.WriteByte(':')

		value, _ := json.Marshal(renderedSegment)
		if segment.mode == fieldMode && segment.value == "extra" {
			value = []byte(renderedSegment)
		}
		renderedRecord.Write(value)
		written++
	}

	renderedRecord.WriteString("}\n")

	result := make([]byte, renderedRecord.Len())
	copy(result, renderedRecord.Bytes())
	return result
}
