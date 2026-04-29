package solislog_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/DasKaroWow/solislog"
)

func TestLoggerWritesMessage(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel,
			&solislog.HandlerOptions{
				Template: "{level} | {message}\n",
			},
		),
	)

	logger.Info("hello")

	want := "INFO | hello\n"
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}

func TestLoggerWritesFullExtraInTemplate(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		solislog.Extra{
			"source": "telegram",
			"id":     "123",
		},
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {message} | {extra}\n",
		}),
	)

	logger.Info("hello")

	got := buf.String()

	if !strings.Contains(got, "INFO | hello | ") {
		t.Fatalf("output = %q, want level and message", got)
	}

	if !strings.Contains(got, `"source":"telegram"`) {
		t.Fatalf("output = %q, want source in extra JSON", got)
	}

	if !strings.Contains(got, `"id":"123"`) {
		t.Fatalf("output = %q, want id in extra JSON", got)
	}
}

func TestNewHandlerWithNilOptionsUsesDefaultTemplate(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel, nil),
	)

	logger.Info("hello")

	got := buf.String()

	if !strings.Contains(got, "INFO") {
		t.Fatalf("output = %q, want level", got)
	}

	if !strings.Contains(got, "hello") {
		t.Fatalf("output = %q, want message", got)
	}
}
