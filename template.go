package solislog

import (
	"fmt"
	"strings"
)

type segmentMode int

const (
	fieldMode segmentMode = iota
	extraMode
	textMode
)

type templateSegment struct {
	mode  segmentMode
	value string
	color string
}

func checkPlaceholderAvailable(placeholderName string) bool {
	switch placeholderName {
	case "time", "level", "message", "extra":
		return true
	default:
		return false
	}
}

func checkColorAvailable(colorName string) bool {
	_, ok := ansiColors[colorName]
	return ok || colorName == "level"
}

func parsePlaceholder(value string) (string, segmentMode) {
	if !strings.HasPrefix(value, "extra[") || !strings.HasSuffix(value, "]") {
		return value, fieldMode
	}

	key := value[6 : len(value)-1]
	return key, extraMode
}

func parseTokens(tokens []templateToken) []templateSegment {
	segments := make([]templateSegment, 0, len(tokens))
	colorStack := newStack[string](len(tokens))

	for _, token := range tokens {
		currentColor, _ := colorStack.peek()

		switch token.kind {
		case tokenColorOpen:
			if !checkColorAvailable(token.value) {
				panic(fmt.Sprintf("unknown color \"%s\"", token.value))
			}
			colorStack.push(token.value)

		case tokenColorClose:
			if !checkColorAvailable(token.value) {
				panic(fmt.Sprintf("unknown color \"%s\"", token.value))
			}
			previousColor, ok := colorStack.pop()
			if !ok || previousColor != token.value {
				panic(fmt.Sprintf("unmatched closing tag \"%s\"", token.value))
			}

		case tokenText:
			segments = append(segments, templateSegment{textMode, token.value, currentColor})

		case tokenPlaceholder:
			value, mode := parsePlaceholder(token.value)
			switch mode {
			case extraMode:
				if value == "" {
					panic("empty extra placeholder")
				}
			case fieldMode:
				if !checkPlaceholderAvailable(value) {
					panic(fmt.Sprintf("unknown placeholder \"%s\"", value))
				}
			}
			segments = append(segments, templateSegment{mode, value, currentColor})
		}
	}

	if colorStack.len() != 0 {
		panic("unclosed color tags remain")
	}

	return segments
}
