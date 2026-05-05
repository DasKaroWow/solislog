package solislog

import "testing"

func TestTokenizeParsesTextColorAndPlaceholder(t *testing.T) {
	tokens := tokenize("<red>abc</red>abc{level}asb")

	want := []templateToken{
		{kind: tokenColorOpen, value: "red"},
		{kind: tokenText, value: "abc"},
		{kind: tokenColorClose, value: "red"},
		{kind: tokenText, value: "abc"},
		{kind: tokenPlaceholder, value: "level"},
		{kind: tokenText, value: "asb"},
	}

	assertTokensEqual(t, tokens, want)
}

func TestTokenizeEscapesSpecialCharacters(t *testing.T) {
	tokens := tokenize(`\<red\>\{level\}\<\/red\>`)

	want := []templateToken{
		{kind: tokenText, value: "<red>{level}</red>"},
	}

	assertTokensEqual(t, tokens, want)
}

func TestTokenizePanicsOnUnclosedPlaceholder(t *testing.T) {
	assertPanics(t, func() {
		tokenize("{level")
	})
}

func TestTokenizePanicsOnUnexpectedClosingBrace(t *testing.T) {
	assertPanics(t, func() {
		tokenize("level}")
	})
}

func TestTokenizePanicsOnUnclosedColorTag(t *testing.T) {
	assertPanics(t, func() {
		tokenize("<red")
	})
}

func TestTokenizePanicsOnUnexpectedClosingAngleBracket(t *testing.T) {
	assertPanics(t, func() {
		tokenize("red>")
	})
}

func TestTokenizePanicsOnEmptyColorTag(t *testing.T) {
	assertPanics(t, func() {
		tokenize("<>")
	})
}

func TestTokenizePanicsOnDanglingEscape(t *testing.T) {
	assertPanics(t, func() {
		tokenize(`abc\`)
	})
}

func assertTokensEqual(t *testing.T, got []templateToken, want []templateToken) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(tokens) = %d, want %d\ngot:  %+v\nwant: %+v", len(got), len(want), got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("tokens[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}
