package solislog

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerWithNilOptionsUsesDefaults(t *testing.T) {
	var buf bytes.Buffer

	handler := NewHandler(&buf, InfoLevel, nil)

	assert.Same(t, &buf, handler.out, "handler.out was not set")
	assert.Equal(t, InfoLevel, handler.level, "handler.level mismatch")
	assert.Equal(t, time.RFC3339, handler.options.TimeFormat, "handler.timeFormat mismatch")
	assert.Equal(t, time.Local, handler.options.Location, "handler.location mismatch")
	assert.NotEmpty(t, handler.template, "handler.template should not be empty")
}

func TestNewHandlerUsesOptions(t *testing.T) {
	var buf bytes.Buffer
	wantTimeFormat := time.DateTime
	wantLocation := time.UTC

	handler := NewHandler(&buf, WarningLevel, &HandlerOptions{
		Template:   "{level} | {message}\n",
		TimeFormat: wantTimeFormat,
		Location:   wantLocation,
		JSON:       true,
	})

	assert.Equal(t, WarningLevel, handler.level, "handler.level mismatch")
	assert.Equal(t, time.DateTime, handler.options.TimeFormat, "handler.timeFormat mismatch")
	assert.Equal(t, wantLocation, handler.options.Location, "handler.location mismatch")
	assert.Equal(t, true, handler.options.JSON, "handler.json should be true")
	assert.NotEmpty(t, handler.template, "handler.template should not be empty")
}
