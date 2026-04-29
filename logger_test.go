package solislog_test

import (
	"bytes"
	"testing"

	"github.com/DasKaroWow/solislog"
)

func TestLoggerWritesMessage(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel, "{level} | {message}\n"),
	)

	logger.Info("hello")

	want := "INFO | hello\n"
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}
