package solislog_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/DasKaroWow/solislog"
)

func TestLoggerWritesJSONMessage(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		solislog.Extra{
			"source": "telegram",
			"id":     "-1",
		},
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template:   "{time} {level} {message} {extra[id]} {extra}",
			JSON:       true,
			TimeFormat: time.RFC3339,
			Location:   time.UTC,
		}),
	)

	logger.Info("hello")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}

	if got["level"] != "INFO" {
		t.Fatalf("level = %v, want %q", got["level"], "INFO")
	}

	if got["message"] != "hello" {
		t.Fatalf("message = %v, want %q", got["message"], "hello")
	}

	if got["id"] != "-1" {
		t.Fatalf("id = %v, want %q", got["id"], "-1")
	}

	extra, ok := got["extra"].(map[string]any)
	if !ok {
		t.Fatalf("extra = %T, want object", got["extra"])
	}

	if extra["source"] != "telegram" {
		t.Fatalf("extra.source = %v, want %q", extra["source"], "telegram")
	}

	if extra["id"] != "-1" {
		t.Fatalf("extra.id = %v, want %q", extra["id"], "-1")
	}
}

func TestLoggerJSONIgnoresColors(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "<red>{level}</red> <level>{message}</level>",
			JSON:     true,
		}),
	)

	logger.Error("boom")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}

	if got["level"] != "ERROR" {
		t.Fatalf("level = %v, want %q", got["level"], "ERROR")
	}

	if got["message"] != "boom" {
		t.Fatalf("message = %v, want %q", got["message"], "boom")
	}
}
