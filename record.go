package solislog

import (
	"encoding/json"
	"strings"
	"time"
)

type record struct {
	time    time.Time
	level   Level
	extra   Extra
	message string
}

func renderField(handler *Handler, rec *record, part *templatePart) string {
	switch part.value {
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

func renderPart(handler *Handler, rec *record, part *templatePart) string {
	switch part.mode {
	case fieldMode:
		return renderField(handler, rec, part)
	case extraMode:
		return rec.extra[part.value]
	}

	return part.value
}

func renderTemplateRecord(handler *Handler, rec *record) string {
	var renderedRecord strings.Builder

	for i := range handler.template {
		part := &handler.template[i]
		renderedPart := renderPart(handler, rec, part)
		renderedRecord.WriteString(renderedPart)
	}

	return renderedRecord.String()
}

func renderJSONRecord(handler *Handler, rec *record) string {
	var rendered strings.Builder
	rendered.WriteByte('{')

	written := 0

	for i := range handler.template {
		part := &handler.template[i]
		if part.mode == textMode {
			continue
		}

		if written > 0 {
			rendered.WriteByte(',')
		}

		renderedPart := renderPart(handler, rec, part)

		name, _ := json.Marshal(part.value)
		rendered.Write(name)
		rendered.WriteByte(':')

		value, _ := json.Marshal(renderedPart)
		if part.mode == fieldMode && part.value == "extra" {
			value = []byte(renderedPart)
		}
		rendered.Write(value)
		written++
	}

	rendered.WriteString("}\n")
	return rendered.String()
}
