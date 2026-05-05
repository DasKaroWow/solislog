package solislog

import (
	"bytes"
	"strings"
	"testing"
)

func TestBeforeHookCanModifyRecord(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(nil, NewHandler(&buf, InfoLevel, &HandlerOptions{
		Template: "{message} | {extra[hooked]}\n",
		BeforeHook: func(record *Record) {
			record.Message = "changed"
			record.Extra["hooked"] = "yes"
		},
	}))

	logger.Info("original")

	want := "changed | yes\n"
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}

func TestAfterHookReceivesRenderedMessage(t *testing.T) {
	var buf bytes.Buffer

	var gotMessage string
	var gotRendered string

	logger := NewLogger(nil, NewHandler(&buf, InfoLevel, &HandlerOptions{
		Template: "{level} | {message}\n",
		AfterHook: func(record *Record, msg string) {
			gotMessage = record.Message
			gotRendered = msg
		},
	}))

	logger.Info("hello")

	want := "INFO | hello\n"

	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}

	if gotMessage != "hello" {
		t.Fatalf("hook message = %q, want %q", gotMessage, "hello")
	}

	if gotRendered != want {
		t.Fatalf("hook rendered = %q, want %q", gotRendered, want)
	}
}

func TestBeforeHookIsIsolatedPerHandler(t *testing.T) {
	var first bytes.Buffer
	var second bytes.Buffer

	logger := NewLogger(
		nil,
		NewHandler(&first, InfoLevel, &HandlerOptions{
			Template: "{message}\n",
			BeforeHook: func(record *Record) {
				record.Message = "changed"
			},
		}),
		NewHandler(&second, InfoLevel, &HandlerOptions{
			Template: "{message}\n",
		}),
	)

	logger.Info("original")

	if first.String() != "changed\n" {
		t.Fatalf("first output = %q, want %q", first.String(), "changed\n")
	}

	if second.String() != "original\n" {
		t.Fatalf("second output = %q, want %q", second.String(), "original\n")
	}
}

func TestAfterHookRunsAfterUnlock(t *testing.T) {
	var buf bytes.Buffer

	var logger *Logger

	logger = NewLogger(nil, NewHandler(&buf, InfoLevel, &HandlerOptions{
		Template: "{message}\n",
		AfterHook: func(record *Record, msg string) {
			if strings.TrimSpace(msg) == "first" {
				logger.Info("second")
			}
		},
	}))

	logger.Info("first")

	want := "first\nsecond\n"
	if buf.String() != want {
		t.Fatalf("output = %q, want %q", buf.String(), want)
	}
}
