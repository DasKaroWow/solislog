package solislog_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/DasKaroWow/solislog"
)

func TestLoggerCanBeUsedConcurrently(t *testing.T) {
	var buf bytes.Buffer

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {message}\n",
		}),
	)

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			logger.Info("hello")
		}()
	}

	wg.Wait()

	output := buf.String()

	got := strings.Count(output, "INFO | hello\n")
	if got != goroutines {
		t.Fatalf("logged lines = %d, want %d\noutput:\n%s", got, goroutines, output)
	}
}

func TestBoundLoggersCanBeUsedConcurrently(t *testing.T) {
	var buf bytes.Buffer

	base := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {extra[id]} | {message}\n",
		}),
	)

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		logger := base.Bind(solislog.Extra{
			"id": "worker",
		})

		go func(logger *solislog.Logger) {
			defer wg.Done()
			logger.Info("hello")
		}(logger)
	}

	wg.Wait()

	output := buf.String()

	got := strings.Count(output, "INFO | worker | hello\n")
	if got != goroutines {
		t.Fatalf("logged lines = %d, want %d\noutput:\n%s", got, goroutines, output)
	}
}
