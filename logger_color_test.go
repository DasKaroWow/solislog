package solislog_test

import (
	"bytes"
	"testing"

	"github.com/DasKaroWow/solislog"
)

func TestLoggerWritesColoredLevel(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "<red>{level}</red> | {message}\n",
		}),
	)

	logger.Info("hello")

	want := "\x1b[31mINFO\x1b[0m | hello\n"
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}

func TestLoggerWritesLevelColoredByLevel(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "<level>{level}</level> | {message}\n",
		}),
	)

	logger.Warning("careful")

	want := "\x1b[33mWARNING\x1b[0m | careful\n"
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}
