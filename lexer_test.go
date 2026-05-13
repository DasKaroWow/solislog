package solislog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func requireTokensEqual(t *testing.T, got, want []templateToken) {
	t.Helper()
	assert.Equal(t, len(want), len(got), "token count mismatch")
	for i := range want {
		assert.Equal(t, want[i], got[i], "token[%d] mismatch", i)
	}
}

func TestScanTemplateTextOnly(t *testing.T) {
	tokens := scanTemplate("hello world")
	requireTokensEqual(t, tokens, []templateToken{
		{kind: tokenText, value: "hello world"},
	})
}

func TestScanTemplatePlaceholdersOnly(t *testing.T) {
	tokens := scanTemplate("{time}{level}")
	requireTokensEqual(t, tokens, []templateToken{
		{kind: tokenPlaceholder, value: "time"},
		{kind: tokenPlaceholder, value: "level"},
	})
}

func TestScanTemplateColorTagsOnly(t *testing.T) {
	tokens := scanTemplate("<red>error</red>")
	requireTokensEqual(t, tokens, []templateToken{
		{kind: tokenColorOpen, value: "red"},
		{kind: tokenText, value: "error"},
		{kind: tokenColorClose, value: "red"},
	})
}

func TestScanTemplateMixedContent(t *testing.T) {
	tokens := scanTemplate("<cyan>{level}</cyan> | {message}")
	requireTokensEqual(t, tokens, []templateToken{
		{kind: tokenColorOpen, value: "cyan"},
		{kind: tokenPlaceholder, value: "level"},
		{kind: tokenColorClose, value: "cyan"},
		{kind: tokenText, value: " | "},
		{kind: tokenPlaceholder, value: "message"},
	})
}

func TestScanTemplateExtraPlaceholder(t *testing.T) {
	tokens := scanTemplate("{extra[request_id]}")
	requireTokensEqual(t, tokens, []templateToken{
		{kind: tokenPlaceholder, value: "extra[request_id]"},
	})
}

func TestScanTemplateEmptyInput(t *testing.T) {
	tokens := scanTemplate("")
	assert.Empty(t, tokens, "empty template should produce zero tokens")
}

func TestScanTemplateEscapesSpecialCharacters(t *testing.T) {
	tokens := scanTemplate(`literal \{placeholder\} and \<color\> and backslash \\`)
	requireTokensEqual(t, tokens, []templateToken{
		{kind: tokenText, value: "literal {placeholder} and <color> and backslash \\"},
	})
}

func TestScanTemplateFlushesConsecutiveText(t *testing.T) {
	tokens := scanTemplate("a<red>b</red>c")
	requireTokensEqual(t, tokens, []templateToken{
		{kind: tokenText, value: "a"},
		{kind: tokenColorOpen, value: "red"},
		{kind: tokenText, value: "b"},
		{kind: tokenColorClose, value: "red"},
		{kind: tokenText, value: "c"},
	})
}

func TestScanTemplatePanicsOnUnclosedPlaceholder(t *testing.T) {
	assert.Panics(t, func() { scanTemplate("{level") })
	assert.Panics(t, func() { scanTemplate("level}") })
}

func TestScanTemplatePanicsOnUnexpectedClosingAngleBracket(t *testing.T) {
	assert.Panics(t, func() { scanTemplate("<text") })
	assert.Panics(t, func() { scanTemplate("text>") })
}

func TestScanTemplatePanicsOnDanglingEscape(t *testing.T) {
	assert.Panics(t, func() { scanTemplate(`text\`) })
}
