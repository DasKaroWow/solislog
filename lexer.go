package solislog

import (
	"strings"
)

type tokenKind int

const (
	tokenText tokenKind = iota
	tokenColorOpen
	tokenColorClose
	tokenPlaceholder
)

type templateToken struct {
	kind  tokenKind
	value string
}

func tokenize(rawTemplate string) []templateToken {
	tokens := make([]templateToken, 0)
	buffer := strings.Builder{}

	flushText := func() {
		if buffer.Len() == 0 {
			return
		}
		tokens = append(tokens, templateToken{tokenText, buffer.String()})
		buffer.Reset()
	}

	for i := 0; i < len(rawTemplate); i++ {
		switch rawTemplate[i] {
		case '\\':
			if i+1 == len(rawTemplate) {
				panic("dangling escape at end of template")
			}
			i++
			buffer.WriteByte(rawTemplate[i])

		case '<':
			flushText()

			endRelativeIndex := strings.IndexByte(rawTemplate[i+1:], '>')
			if endRelativeIndex == -1 {
				panic("unclosed color tag")
			}

			endAbsoluteIndex := i + endRelativeIndex + 1
			tag := rawTemplate[i+1 : endAbsoluteIndex]
			if tag == "" || tag == "/" {
				panic("empty color tag")
			}

			if strings.HasPrefix(tag, "/") {
				tokens = append(tokens, templateToken{tokenColorClose, tag[1:]})
			} else {
				tokens = append(tokens, templateToken{tokenColorOpen, tag})
			}
			i = endAbsoluteIndex

		case '{':
			flushText()

			endRelativeIndex := strings.IndexByte(rawTemplate[i+1:], '}')
			if endRelativeIndex == -1 {
				panic("unclosed placeholder")
			}

			endAbsoluteIndex := i + endRelativeIndex + 1
			placeholder := rawTemplate[i+1 : endAbsoluteIndex]
			if placeholder == "" {
				panic("empty placeholder")
			}
			tokens = append(tokens, templateToken{tokenPlaceholder, placeholder})
			i = endAbsoluteIndex

		case '}':
			panic("unexpected closing brace")

		case '>':
			panic("unexpected closing angle bracket")

		default:
			buffer.WriteByte(rawTemplate[i])
		}
	}

	flushText()
	return tokens
}
