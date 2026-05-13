package solislog_test

import (
	"io"
	"sync"
	"testing"
	"time"

	"github.com/DasKaroWow/solislog"
)

func BenchmarkLoggerInfoText(b *testing.B) {
	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(io.Discard, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {message}\n",
		}),
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello")
	}
}

func BenchmarkLoggerInfoWithExtra(b *testing.B) {
	logger := solislog.NewLogger(
		solislog.Extra{
			"service": "api",
			"env":     "dev",
			"id":      "123",
		},
		solislog.NewHandler(io.Discard, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {extra[service]} | {extra[env]} | {extra[id]} | {message}\n",
		}),
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello")
	}
}

func BenchmarkLoggerInfoJSON(b *testing.B) {
	logger := solislog.NewLogger(
		solislog.Extra{
			"service": "api",
			"env":     "dev",
			"id":      "123",
		},
		solislog.NewHandler(io.Discard, solislog.InfoLevel, &solislog.HandlerOptions{
			JSON:       true,
			TimeFormat: time.RFC3339,
			Location:   time.UTC,
			Template:   "{time} {level} {message} {extra[id]} {extra}",
		}),
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello")
	}
}

func BenchmarkLoggerInfoFilteredOut(b *testing.B) {
	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(io.Discard, solislog.ErrorLevel, &solislog.HandlerOptions{
			Template: "{level} | {message}\n",
		}),
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello")
	}
}

func BenchmarkLoggerInfoParallel(b *testing.B) {
	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(io.Discard, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {message}\n",
		}),
	)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("hello")
		}
	})
}

func BenchmarkBoundLoggerInfo(b *testing.B) {
	base := solislog.NewLogger(
		solislog.Extra{
			"service": "api",
		},
		solislog.NewHandler(io.Discard, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {extra[service]} | {extra[request_id]} | {message}\n",
		}),
	)

	logger := base.Bind(solislog.Extra{
		"request_id": "req-123",
	})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello")
	}
}

func BenchmarkLoggerInfoMultipleHandlers(b *testing.B) {
	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(io.Discard, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {message}\n",
		}),
		solislog.NewHandler(io.Discard, solislog.InfoLevel, &solislog.HandlerOptions{
			JSON:     true,
			Template: "{level} {message}",
		}),
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello")
	}
}

type lockedDiscard struct {
	mu sync.Mutex
}

func (w *lockedDiscard) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(p), nil
}

func BenchmarkLoggerInfoLockedWriter(b *testing.B) {
	writer := &lockedDiscard{}

	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(writer, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "{level} | {message}\n",
		}),
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("hello")
	}
}
