package solislog_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/DasKaroWow/solislog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Log(buf.String())

	var got map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &got), "output is not valid JSON")

	level := got["level"]
	message := got["message"]
	id := got["id"]

	assert.Equal(t, "INFO", level)
	assert.Equal(t, "hello", message)
	assert.Equal(t, "-1", id)

	extraAny := got["extra"]
	extra, ok := extraAny.(map[string]any)
	require.True(t, ok, "extra field type mismatch: got %T, want object", extraAny)

	assert.Equal(t, "telegram", extra["source"])
	assert.Equal(t, "-1", extra["id"])
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

	t.Log(buf.String())

	var got map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &got), "output is not valid JSON")

	level := got["level"]
	message := got["message"]

	assert.Equal(t, "ERROR", level)
	assert.Equal(t, "boom", message)
}
