package solislog

import (
	"strings"
	"time"
)

type templatePart struct {
	isField bool
	key     string
	value   string
}

var availableParts = map[string]struct{}{
	"time":    {},
	"level":   {},
	"message": {},
}

func (part *templatePart) checkExtra() {
	if !part.isField {
		return
	}

	if strings.HasPrefix(part.value, "extra[") && strings.HasSuffix(part.value, "]") {
		part.key = part.value[6 : len(part.value)-1]
	} else {
		_, found := availableParts[part.value]
		if !found {
			panic("Unknown field")
		}
	}
}

func renderRecord(parts []templatePart, rec *record) string {
	var renderedRecord strings.Builder

	for _, part := range parts {
		rendered := renderField(part, rec)
		renderedRecord.WriteString(rendered)
	}
	return renderedRecord.String()
}

func renderField(part templatePart, rec *record) string {
	if !part.isField {
		return part.value
	}

	switch part.value {
	case "time":
		return rec.time.Format(time.RFC3339)
	case "level":
		return rec.level.String()
	case "message":
		return rec.message
	default:
		return rec.extra[part.key]
	}
}

func parseTemplate(rawTemplate string) []templatePart {
	if rawTemplate == "" {
		return []templatePart{
			{
				isField: true,
				value:   "time",
			},
			{
				isField: false,
				value:   " | ",
			},
			{
				isField: true,
				value:   "level",
			},
			{
				isField: false,
				value:   " | ",
			},
			{
				isField: true,
				value:   "message",
			},
			{
				isField: false,
				value:   "\n",
			},
		}
	}
	var buf strings.Builder
	template := make([]templatePart, 0, len(rawTemplate)/10)

	for i := 0; i < len(rawTemplate); i++ {
		switch rawTemplate[i] {
		case '{':
			if buf.Len() > 0 {
				template = append(template, templatePart{
					isField: false,
					value:   buf.String(),
				})
				buf.Reset()
			}

			j := i + 1
			for j < len(rawTemplate) && rawTemplate[j] != '}' {
				j++
			}
			if j == len(rawTemplate) {
				panic("unclosed placeholder")
			}

			field := rawTemplate[i+1 : j]
			if field == "" {
				panic("empty placeholder")
			}

			template = append(template, templatePart{
				isField: true,
				value:   field,
			})
			template[len(template)-1].checkExtra()
			i = j

		case '}':
			panic("unexpected closing brace")

		default:
			buf.WriteByte(rawTemplate[i])
		}
	}

	if buf.Len() > 0 {
		template = append(template, templatePart{
			isField: false,
			value:   buf.String(),
		})
	}

	return template
}
