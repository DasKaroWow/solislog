package solislog_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/DasKaroWow/solislog"
	"github.com/stretchr/testify/require"
)

type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (w *lockedBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

func (w *lockedBuffer) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}
func TestLoggerCanBeUsedConcurrently(t *testing.T) {
	var buf lockedBuffer // ← теперь потокобезопасный
	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(&buf, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {message}\n",
		}),
	)

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			logger.Info("hello")
		}()
	}

	wg.Wait()

	output := buf.String()
	t.Log(output)

	want := goroutines
	got := strings.Count(output, "INFO | hello\n")

	require.Equal(t, want, got, "logged lines count mismatch")
}

func TestBoundLoggersCanBeUsedConcurrently(t *testing.T) {
	var buf lockedBuffer
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
		boundLogger := base.Bind(solislog.Extra{"id": "worker"})
		go func(l *solislog.Logger) {
			defer wg.Done()
			l.Info("hello")
		}(boundLogger)
	}

	wg.Wait()

	output := buf.String()
	t.Log(output)

	want := goroutines
	got := strings.Count(output, "INFO | worker | hello\n")

	require.Equal(t, want, got, "bound logger logged lines count mismatch")
}
