package solislog_test

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/DasKaroWow/solislog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggerWritesCallerMetadata(t *testing.T) {
	var buf bytes.Buffer
	logger := solislog.NewLogger(nil,
		solislog.NewHandler(&buf, solislog.InfoLevel,
			&solislog.HandlerOptions{
				Template:   "{file}|{path}|{line}|{function}|{caller}|{message}\n",
				WithCaller: true,
			},
		),
	)

	logger.Info("test")

	output := strings.TrimSpace(buf.String())
	parts := strings.Split(output, "|")

	require.Len(t, parts, 6)

	t.Log(output)

	file := parts[0]
	_ = parts[1]
	line := parts[2]
	function := parts[3]
	caller := parts[4]
	message := parts[5]

	assert.Equal(t, "logger_metadata_test.go", filepath.Base(file))
	lineNumber, err := strconv.Atoi(line)
	assert.NoError(t, err)
	assert.Greater(t, lineNumber, 0)
	assert.Contains(t, function, ".TestLoggerWritesCallerMetadata")
	assert.Contains(t, caller, "logger_metadata_test.go:")
	assert.Equal(t, "test", message)
}

func TestLoggerWritesCallerMetadataInJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := solislog.NewLogger(nil,
		solislog.NewHandler(&buf, solislog.InfoLevel,
			&solislog.HandlerOptions{
				JSON:       true,
				WithCaller: true,
				Template:   "{file} {path} {line} {function} {caller} {message}",
			},
		),
	)

	logger.Info("hello metadata")

	var got map[string]string
	require.NoError(t, json.Unmarshal(buf.Bytes(), &got), "output is not valid JSON: %s", buf.String())

	t.Log(buf.String())

	assert.Equal(t, "logger_metadata_test.go", got["file"])
	assert.Equal(t, "logger_metadata_test.go", filepath.Base(got["path"]))
	assert.NotEmpty(t, got["line"])
	assert.True(t, strings.HasSuffix(got["function"], ".TestLoggerWritesCallerMetadataInJSON"))
	assert.True(t, strings.HasPrefix(got["caller"], "logger_metadata_test.go:"))
	assert.Equal(t, "hello metadata", got["message"])
}
