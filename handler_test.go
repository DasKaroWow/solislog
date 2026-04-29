package solislog

import (
	"bytes"
	"testing"
	"time"
)

func TestNewHandlerWithNilOptionsUsesDefaults(t *testing.T) {
	var buf bytes.Buffer

	handler := NewHandler(&buf, InfoLevel, nil)

	if handler.out != &buf {
		t.Fatal("handler.out was not set")
	}

	if handler.level != InfoLevel {
		t.Fatalf("handler.level = %v, want %v", handler.level, InfoLevel)
	}

	if handler.timeFormat != time.RFC3339 {
		t.Fatalf("handler.timeFormat = %q, want %q", handler.timeFormat, time.RFC3339)
	}

	if handler.location != time.Local {
		t.Fatalf("handler.location = %v, want %v", handler.location, time.Local)
	}

	if handler.json {
		t.Fatal("handler.json = true, want false")
	}

	if len(handler.template) == 0 {
		t.Fatal("handler.template is empty")
	}
}

func TestNewHandlerUsesOptions(t *testing.T) {
	var buf bytes.Buffer
	location := time.UTC

	handler := NewHandler(&buf, WarningLevel, &HandlerOptions{
		Template:   "{level} | {message}\n",
		TimeFormat: time.DateTime,
		Location:   location,
		JSON:       true,
	})

	if handler.level != WarningLevel {
		t.Fatalf("handler.level = %v, want %v", handler.level, WarningLevel)
	}

	if handler.timeFormat != time.DateTime {
		t.Fatalf("handler.timeFormat = %q, want %q", handler.timeFormat, time.DateTime)
	}

	if handler.location != location {
		t.Fatalf("handler.location = %v, want %v", handler.location, location)
	}

	if !handler.json {
		t.Fatal("handler.json = false, want true")
	}

	if len(handler.template) == 0 {
		t.Fatal("handler.template is empty")
	}
}
