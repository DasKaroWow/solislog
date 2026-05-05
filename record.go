package solislog

import (
	"encoding/json"
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

type record struct {
	time    time.Time
	level   Level
	extra   Extra
	message string
}

func renderField(handler *Handler, rec *record, segment *templateSegment) string {
	switch segment.value {
	case "time":
		return rec.time.In(handler.location).Format(handler.timeFormat)
	case "level":
		return rec.level.String()
	case "message":
		return rec.message
	case "extra":
		data, err := json.Marshal(rec.extra)
		if err != nil {
			return "{}"
		}
		return string(data)
	}
	return ""
}

func renderSegment(handler *Handler, rec *record, segment *templateSegment) string {
	switch segment.mode {
	case fieldMode:
		return renderField(handler, rec, segment)
	case extraMode:
		return rec.extra[segment.value]
	}

	return segment.value
}

func renderTemplateRecord(handler *Handler, rec *record) string {
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
			renderedRecord.WriteString(renderColor(segment.color, rec.level))
			previousColor = segment.color
		}
		renderedRecord.WriteString(renderSegment(handler, rec, segment))
	}

	if previousColor != "" {
		renderedRecord.WriteString(ansiReset)
	}
	return renderedRecord.String()
}

func renderJSONRecord(handler *Handler, rec *record) string {
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
