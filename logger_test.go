package solislog_test

import (
	"bytes"
	"testing"

	"github.com/DasKaroWow/solislog"
	"github.com/stretchr/testify/assert"
)

func TestLoggerWritesMessage(t *testing.T) {
	var buffer bytes.Buffer

	logger := solislog.NewLogger(
		solislog.Extra{"id": "1"},
		solislog.NewHandler(&buffer, solislog.InfoLevel,
			&solislog.HandlerOptions{
				Template: "{level} | {extra[id]} | {message}\n",
			},
		),
	)

	logger.Info("hello")
	want := "INFO | 1 | hello\n"
	got := buffer.String()
	assert.Equal(t, got, want)
}

func TestLoggerWritesFullExtraInTemplate(t *testing.T) {
	var buffer bytes.Buffer

	logger := solislog.NewLogger(
		solislog.Extra{
			"source": "unknown",
			"id":     "-1",
		},
		solislog.NewHandler(&buffer, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {extra} | {message}\n",
		}),
	)

	logger.Info("hello")
	// want := `INFO | {"source": "unknown", "id": "-1"} | hello` + "\n"
	got := buffer.String()

	assert.Contains(t, got, `"source":"unknown"`)
	assert.Contains(t, got, `"id":"-1"`)
	assert.Contains(t, got, "INFO |")
	assert.Contains(t, got, "| hello\n")
}

func TestNewHandlerWithNilOptionsUsesDefaultTemplate(t *testing.T) {
	var buffer bytes.Buffer

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buffer, solislog.InfoLevel, nil),
	)

	logger.Info("hello")
	got := buffer.String()
	assert.Contains(t, got, "| INFO | hello\n")
}

func TestBindExtraOverridesTemplate(t *testing.T) {
	var buf bytes.Buffer
	base := solislog.NewLogger(
		solislog.Extra{"service": "api"},
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {extra[service]} | {extra[request_id]} | {message}\n",
		}),
	)

	bound := base.Bind(solislog.Extra{"request_id": "req-1", "service": "web"})
	bound.Info("hello")

	assert.Equal(t, buf.String(), "INFO | web | req-1 | hello\n")
}
