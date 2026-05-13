package solislog_test

import (
	"bytes"
	"testing"

	"github.com/DasKaroWow/solislog"
	"github.com/stretchr/testify/assert"
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
	got := buf.String()
	assert.Equal(t, want, got)
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
	got := buf.String()
	assert.Equal(t, want, got)
}

func TestLoggerWritesNestingColors(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "<red>a<yellow>b</yellow>c</red> | {message}\n",
		}),
	)

	logger.Info("test")

	want := "\x1b[31ma\x1b[33mb\x1b[31mc\x1b[0m | test\n"
	got := buf.String()
	assert.Equal(t, want, got)
}
