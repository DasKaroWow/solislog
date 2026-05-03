package solislog

import (
	"fmt"
	"strings"
)

type coloredSegment struct {
	text  string
	color string
}

var availableColors = map[string]struct{}{
	"red":   {},
	"green": {},
	"blue":  {},
	"white": {},
}

// func findTokenEnd(template string, tokenSymbol byte, offset int) int {
// 	index := strings.IndexByte(template[offset:], tokenSymbol)
// 	if index == -1 {
// 		panic("unclosed placeholder")
// 	}
// 	return index + offset
// }

func parseColors(rawTemplate string) []coloredSegment {
	var buffer strings.Builder
	colorStack := newStack[string](len(rawTemplate) / 10)
	segments := make([]coloredSegment, 0, len(rawTemplate)/10)

	flush := func() {
		if buffer.Len() == 0 {
			return
		}
		segments = append(segments, coloredSegment{
			text:  buffer.String(),
			color: colorStack.peek(),
		})
		buffer.Reset()
	}

	for i := 0; i < len(rawTemplate); i++ {
		switch rawTemplate[i] {
		case '\\':
			if i == len(rawTemplate)-1 {
				panic("dangling escape character at end of template")
			}
			i++
			buffer.WriteByte(rawTemplate[i])
		case '<':
			end := strings.IndexByte(rawTemplate[i+1:], '>')
			if end == -1 {
				panic("unclosed color tag")
			}

			end = i + 1 + end
			tag := raw[i+1 : end]

			if !isValidColorTag(tag) {
				buf.WriteString(raw[i : end+1]) // не тег, обычный текст
				i = end
				continue
			}

			flush()

			if strings.HasPrefix(tag, "/") {
				fmt.Printf("CLOSE %q\n", tag[1:])
			} else {
				fmt.Printf("OPEN %q\n", tag)
			}

			i = end

		default:
			buf.WriteByte(raw[i])
		}
	}

	flush()
}

func isValidColorTag(tag string) bool {
	if tag == "" {
		return false
	}

	if strings.HasPrefix(tag, "/") {
		tag = tag[1:]
		if tag == "" {
			return false
		}
	}

	for i := 0; i < len(tag); i++ {
		c := tag[i]
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_' ||
			c == '-') {
			return false
		}
	}

	return true
}

func main() {
	parseColors(`ttt<red>abc</red>`)
}
