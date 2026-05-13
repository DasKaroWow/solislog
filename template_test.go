package solislog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireSegmentsEqual(t *testing.T, got, want []templateSegment) {
	t.Helper()
	assert.Equal(t, len(want), len(got), "got: %+v\nwant: %+v", got, want)
	for i := range want {
		require.Equal(t, want[i], got[i])
	}
}

func TestBuildSegmentsEmptyInput(t *testing.T) {
	got := buildSegments(nil)
	assert.Equal(t, 0, len(got))
}

func TestBuildSegmentsTextAndFields(t *testing.T) {
	tokens := []templateToken{
		{kind: tokenPlaceholder, value: "time"},
		{kind: tokenText, value: " | "},
		{kind: tokenPlaceholder, value: "level"},
		{kind: tokenText, value: " | "},
		{kind: tokenPlaceholder, value: "message"},
	}
	want := []templateSegment{
		{mode: fieldMode, value: "time"},
		{mode: textMode, value: " | "},
		{mode: fieldMode, value: "level"},
		{mode: textMode, value: " | "},
		{mode: fieldMode, value: "message"},
	}
	requireSegmentsEqual(t, buildSegments(tokens), want)
}

func TestBuildSegmentsExtraPlaceholders(t *testing.T) {
	tokens := []templateToken{
		{kind: tokenPlaceholder, value: "extra[request_id]"},
		{kind: tokenPlaceholder, value: "extra"},
	}
	want := []templateSegment{
		{mode: extraMode, value: "request_id"},
		{mode: fieldMode, value: "extra"},
	}
	requireSegmentsEqual(t, buildSegments(tokens), want)
}

func TestBuildSegmentsColorsApplyToContent(t *testing.T) {
	tokens := []templateToken{
		{kind: tokenColorOpen, value: "red"},
		{kind: tokenPlaceholder, value: "level"},
		{kind: tokenText, value: " error: "},
		{kind: tokenPlaceholder, value: "message"},
		{kind: tokenColorClose, value: "red"},
	}
	want := []templateSegment{
		{mode: fieldMode, value: "level", color: "red"},
		{mode: textMode, value: " error: ", color: "red"},
		{mode: fieldMode, value: "message", color: "red"},
	}
	requireSegmentsEqual(t, buildSegments(tokens), want)
}

func TestBuildSegmentsNestedColors(t *testing.T) {
	tokens := []templateToken{
		{kind: tokenColorOpen, value: "red"},
		{kind: tokenText, value: "a"},
		{kind: tokenColorOpen, value: "blue"},
		{kind: tokenText, value: "b"},
		{kind: tokenColorClose, value: "blue"},
		{kind: tokenText, value: "c"},
		{kind: tokenColorClose, value: "red"},
	}
	want := []templateSegment{
		{mode: textMode, value: "a", color: "red"},
		{mode: textMode, value: "b", color: "blue"},
		{mode: textMode, value: "c", color: "red"},
	}
	requireSegmentsEqual(t, buildSegments(tokens), want)
}

func TestBuildSegmentsLevelColorSpecialCase(t *testing.T) {
	tokens := []templateToken{
		{kind: tokenColorOpen, value: "level"},
		{kind: tokenPlaceholder, value: "level"},
		{kind: tokenColorClose, value: "level"},
	}
	want := []templateSegment{
		{mode: fieldMode, value: "level", color: "level"},
	}
	requireSegmentsEqual(t, buildSegments(tokens), want)
}

func TestBuildSegmentsPanicsOnUnknownField(t *testing.T) {
	assert.Panics(t, func() {
		buildSegments([]templateToken{
			{kind: tokenPlaceholder, value: "unknown"},
		})
	})
}

func TestBuildSegmentsPanicsOnUnknownColor(t *testing.T) {
	assert.Panics(t, func() {
		buildSegments([]templateToken{
			{kind: tokenColorOpen, value: "neon"},
		})
	})
	assert.Panics(t, func() {
		buildSegments([]templateToken{
			{kind: tokenColorClose, value: "neon"},
		})
	})
}

func TestBuildSegmentsPanicsOnUnmatchedClosingTag(t *testing.T) {
	assert.Panics(t, func() {
		buildSegments([]templateToken{
			{kind: tokenColorClose, value: "red"},
		})
	})
}

func TestBuildSegmentsPanicsOnMismatchedClosingTag(t *testing.T) {
	assert.Panics(t, func() {
		buildSegments([]templateToken{
			{kind: tokenColorOpen, value: "red"},
			{kind: tokenColorClose, value: "blue"},
		})
	})
}

func TestBuildSegmentsPanicsOnUnclosedColorStack(t *testing.T) {
	assert.Panics(t, func() {
		buildSegments([]templateToken{
			{kind: tokenColorOpen, value: "red"},
			{kind: tokenText, value: "text"},
		})
	})
}

func TestBuildSegmentsPanicsOnEmptyExtraKey(t *testing.T) {
	assert.Panics(t, func() {
		buildSegments([]templateToken{
			{kind: tokenPlaceholder, value: "extra[]"},
		})
	})
}
