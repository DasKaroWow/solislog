package solislog

import (
	"slices"
	"strings"
)

type tokenKind int

const (
	textToken tokenKind = iota
	placeholderToken
	openStyleToken
	closeStyleToken
)

type templateToken struct {
	kind  tokenKind
	value string
}

func token(kind tokenKind, value string) templateToken {
	return templateToken{
		kind:  kind,
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

func tokenizeRawTemplate(raw string) []templateToken {
	var tokens []templateToken
	var buf strings.Builder

	flushText := func() {
		if buf.Len() == 0 {
			return
		}

		tokens = append(tokens, token(textToken, buf.String()))
		buf.Reset()
	}

	for i := 0; i < len(raw); i++ {
		switch raw[i] {
		case '\\':
			i += 1
			if i != len(raw) {
				buf.WriteByte(raw[i])
			}
		case '{':
			flushText()

			end := findPlaceholderEnd(raw, i)

			tokens = append(tokens, templateToken{
				kind:  placeholderToken,
				value: raw[i+1 : end],
			})
			i = end

		case '}':
			panic("unexpected closing brace")

		case '<':
			tag, closing, ok, end := lexStyleTag(raw, i)
			if !ok {
				buf.WriteByte(raw[i])
				continue
			}

			flushText()

			kind := openStyleToken
			if closing {
				kind = closeStyleToken
			}

			tokens = append(tokens, templateToken{
				kind:  kind,
				value: tag,
			})

			i = end

		default:
			buf.WriteByte(raw[i])
		}
	}

	flushText()
	return slices.Clip(tokens)
}
