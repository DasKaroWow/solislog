package solislog_test

import (
	"errors"
	"testing"

	"github.com/DasKaroWow/solislog"
)

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestLoggerCallsErrorHandlerOnWriteError(t *testing.T) {
	var gotErr error
	var gotMsg string

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(failingWriter{}, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {message}\n",
			ErrorHandler: func(err error, msg string) {
				gotErr = err
				gotMsg = msg
			},
		}),
	)

	logger.Info("hello")

	if gotErr == nil {
		t.Fatal("error handler was not called")
	}

	if gotErr.Error() != "write failed" {
		t.Fatalf("error = %q, want %q", gotErr.Error(), "write failed")
	}

	wantMsg := "INFO | hello\n"
	if gotMsg != wantMsg {
		t.Fatalf("msg = %q, want %q", gotMsg, wantMsg)
	}
}
