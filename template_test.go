package solislog

import "testing"

func TestParseTemplateParsesTextAndFields(t *testing.T) {
	segments := parseTokens(tokenize("{time} | {level} | {message}\n"))

	want := []templateSegment{
		{mode: fieldMode, value: "time"},
		{mode: textMode, value: " | "},
		{mode: fieldMode, value: "level"},
		{mode: textMode, value: " | "},
		{mode: fieldMode, value: "message"},
		{mode: textMode, value: "\n"},
	}

	assertSegmentsEqual(t, segments, want)
}

func TestParseTemplateParsesExtraField(t *testing.T) {
	segments := parseTokens(tokenize("{extra[source]}"))

	want := []templateSegment{
		{mode: extraMode, value: "source"},
	}

	assertSegmentsEqual(t, segments, want)
}

func TestParseTemplateParsesFullExtraField(t *testing.T) {
	segments := parseTokens(tokenize("{extra}"))

	want := []templateSegment{
		{mode: fieldMode, value: "extra"},
	}

	assertSegmentsEqual(t, segments, want)
}

func TestParseTemplateAddsColorToTextAndFields(t *testing.T) {
	segments := parseTokens(tokenize("<red>{level} | text</red>"))

	want := []templateSegment{
		{mode: fieldMode, value: "level", color: "red"},
		{mode: textMode, value: " | text", color: "red"},
	}

	assertSegmentsEqual(t, segments, want)
}

func TestParseTemplateSupportsNestedColors(t *testing.T) {
	segments := parseTokens(tokenize("<red>a<blue>b</blue>c</red>"))

	want := []templateSegment{
		{mode: textMode, value: "a", color: "red"},
		{mode: textMode, value: "b", color: "blue"},
		{mode: textMode, value: "c", color: "red"},
	}

	assertSegmentsEqual(t, segments, want)
}

func TestParseTemplateSupportsLevelColor(t *testing.T) {
	segments := parseTokens(tokenize("<level>{level}</level>"))

	want := []templateSegment{
		{mode: fieldMode, value: "level", color: "level"},
	}

	assertSegmentsEqual(t, segments, want)
}

func TestParseTemplatePanicsOnUnknownField(t *testing.T) {
	assertPanics(t, func() {
		parseTokens(tokenize("{unknown}"))
	})
}

func TestParseTemplatePanicsOnEmptyExtraKey(t *testing.T) {
	assertPanics(t, func() {
		parseTokens(tokenize("{extra[]}"))
	})
}

func TestParseTemplatePanicsOnUnknownColor(t *testing.T) {
	assertPanics(t, func() {
		parseTokens(tokenize("<unknown>text</unknown>"))
	})
}

func TestParseTemplatePanicsOnUnmatchedClosingTag(t *testing.T) {
	assertPanics(t, func() {
		parseTokens(tokenize("</red>"))
	})
}

func TestParseTemplatePanicsOnMismatchedClosingTag(t *testing.T) {
	assertPanics(t, func() {
		parseTokens(tokenize("<red>text</blue>"))
	})
}

func TestParseTemplatePanicsOnUnclosedColorTag(t *testing.T) {
	assertPanics(t, func() {
		parseTokens(tokenize("<red>text"))
	})
}

func assertSegmentsEqual(t *testing.T, got []templateSegment, want []templateSegment) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(segments) = %d, want %d\ngot:  %+v\nwant: %+v", len(got), len(want), got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("segments[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()

	fn()
}
