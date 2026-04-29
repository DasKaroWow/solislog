package solislog

import (
	"slices"
	"strings"
)

type partMode int

const (
	fieldMode partMode = iota
	extraMode
	textMode
)

type templatePart struct {
	mode  partMode
	value string
}

var availableParts = map[string]struct{}{
	"time":    {},
	"level":   {},
	"message": {},
	"extra":   {},
}

func parseExtraField(value string) (string, bool) {
	if !strings.HasPrefix(value, "extra[") || !strings.HasSuffix(value, "]") {
		return value, false
	}

	key := value[6 : len(value)-1]
	if key == "" {
		panic("empty extra key")
	}

	return key, true
}

func parsePlaceholder(placeholder string) templatePart {
	if placeholder == "" {
		panic("empty placeholder")
	}

	value, isExtra := parseExtraField(placeholder)
	if isExtra {
		return templatePart{
			mode:  extraMode,
			value: value,
		}
	}

	if _, ok := availableParts[value]; !ok {
		panic("unknown field")
	}

	return templatePart{
		mode:  fieldMode,
		value: value,
	}
}

func findPlaceholderEnd(template string, start int) int {
	for i := start + 1; i < len(template); i++ {
		if template[i] == '}' {
			return i
		}
	}

	panic("unclosed placeholder")
}

func parseTemplate(rawTemplate string) []templatePart {
	var buf strings.Builder
	parts := make([]templatePart, 0, len(rawTemplate)/10)

	for i := 0; i < len(rawTemplate); i++ {
		switch rawTemplate[i] {
		case '{':
			if buf.Len() > 0 {
				parts = append(parts,
					templatePart{
						mode:  textMode,
						value: buf.String(),
					},
				)
				buf.Reset()
			}

			j := findPlaceholderEnd(rawTemplate, i)
			placeholder := rawTemplate[i+1 : j]

			part := parsePlaceholder(placeholder)
			parts = append(parts, part)
			i = j
		case '}':
			panic("unexpected closing brace")
		default:
			buf.WriteByte(rawTemplate[i])
		}
	}

	if buf.Len() > 0 {
		parts = append(parts, templatePart{mode: textMode, value: buf.String()})
	}

	return slices.Clip(parts)
}
