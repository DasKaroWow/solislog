package solislog

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestLoggerWritesCallerMetadata(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(nil, NewHandler(&buf, InfoLevel, &HandlerOptions{
		Template: "{file}|{path}|{line}|{function}|{caller}|{message}\n",
	}))

	logger.Info("hello metadata")

	output := strings.TrimSpace(buf.String())
	parts := strings.Split(output, "|")

	if len(parts) != 6 {
		t.Fatalf("expected 6 parts, got %d: %q", len(parts), output)
	}

	file := parts[0]
	path := parts[1]
	line := parts[2]
	function := parts[3]
	caller := parts[4]
	message := parts[5]

	if file != "logger_metadata_test.go" {
		t.Fatalf("expected file %q, got %q; output: %q", "logger_metadata_test.go", file, output)
	}

	if filepath.Base(path) != "logger_metadata_test.go" {
		t.Fatalf("expected path to point to logger_metadata_test.go, got %q; output: %q", path, output)
	}

	lineNumber, err := strconv.Atoi(line)
	if err != nil {
		t.Fatalf("expected line to be a number, got %q; output: %q", line, output)
	}

	if lineNumber <= 0 {
		t.Fatalf("expected positive line number, got %d; output: %q", lineNumber, output)
	}

	if !strings.HasSuffix(function, ".TestLoggerWritesCallerMetadata") {
		t.Fatalf("expected function to end with .TestLoggerWritesCallerMetadata, got %q; output: %q", function, output)
	}

	expectedCallerPrefix := "logger_metadata_test.go:"
	if !strings.HasPrefix(caller, expectedCallerPrefix) {
		t.Fatalf("expected caller to start with %q, got %q; output: %q", expectedCallerPrefix, caller, output)
	}

	callerLine := strings.TrimPrefix(caller, expectedCallerPrefix)
	if callerLine != line {
		t.Fatalf("expected caller line %q to match line %q; output: %q", callerLine, line, output)
	}

	if message != "hello metadata" {
		t.Fatalf("expected message %q, got %q; output: %q", "hello metadata", message, output)
	}
}

func TestLoggerWritesCallerMetadataInJSON(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(nil, NewHandler(&buf, InfoLevel, &HandlerOptions{
		JSON:     true,
		Template: "{file} {path} {line} {function} {caller} {message}",
	}))

	logger.Info("hello metadata")

	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}

	if got["file"] != "logger_metadata_test.go" {
		t.Fatalf("file = %q, want %q", got["file"], "logger_metadata_test.go")
	}

	if filepath.Base(got["path"]) != "logger_metadata_test.go" {
		t.Fatalf("path = %q, want logger_metadata_test.go", got["path"])
	}

	if got["line"] == "" {
		t.Fatalf("line is empty")
	}

	if !strings.HasSuffix(got["function"], ".TestLoggerWritesCallerMetadataInJSON") {
		t.Fatalf("function = %q, want test function", got["function"])
	}

	if !strings.HasPrefix(got["caller"], "logger_metadata_test.go:") {
		t.Fatalf("caller = %q, want logger_metadata_test.go:<line>", got["caller"])
	}

	if got["message"] != "hello metadata" {
		t.Fatalf("message = %q, want %q", got["message"], "hello metadata")
	}
}
