package solislog

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindReturnsSameInstanceForNilExtra(t *testing.T) {
	base := NewLogger(Extra{"k": "v"}, NewHandler(io.Discard, InfoLevel, nil))

	boundNil := base.Bind(nil)
	assert.Same(t, base, boundNil, "Bind(nil) should return the same logger instance")

	boundEmpty := base.Bind(Extra{})
	assert.NotSame(t, base, boundEmpty, "Bind(Extra{}) should return the same logger instance")
}

func TestBindCreatesNewInstanceWithMergedExtra(t *testing.T) {
	base := NewLogger(Extra{"base": "1"}, NewHandler(io.Discard, InfoLevel, nil))
	bound := base.Bind(Extra{"new": "2"})

	assert.NotSame(t, base, bound, "Bind with non-nil extra should return a new instance")
	assert.Equal(t, bound.extra["base"], "1")
	assert.Equal(t, bound.extra["new"], "2")
}

func TestBindDoesNotMutateOriginalExtra(t *testing.T) {
	base := NewLogger(Extra{"k": "v"}, NewHandler(io.Discard, InfoLevel, nil))
	_ = base.Bind(Extra{"k": "overridden", "new": "val"})

	assert.Equal(t, base.extra["k"], "v", "original extra was mutated")
	if _, exists := base.extra["new"]; exists {
		assert.Fail(t, "original extra received new keys from Bind")
	}
}
