package solislog_test

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/DasKaroWow/solislog"
)

// ExampleNewLogger_text demonstrates a basic text logger with ANSI colors,
// extra fields, and multiple log levels.
func ExampleNewLogger_text() {
	logger := solislog.NewLogger(
		solislog.Extra{
			"service": "api",
			"env":     "dev",
		},
		solislog.NewHandler(os.Stdout, solislog.DebugLevel, &solislog.HandlerOptions{
			Template: "<gray>{time}</gray> | <level>{level}</level> | service={extra[service]} env={extra[env]} | {message}\n",
		}),
	)

	logger.Debug("debug message")
	logger.Info("server started")
	logger.Warning("slow request")
	logger.Error("request failed")

	// NOTE: // Output: is intentionally omitted because the template
	// contains {time} and ANSI escape codes, which change per run.
	// The function will still compile and run during `go test`.
}

// ExampleLogger_Infof demonstrates formatted logging.
func ExampleLogger_Infof() {
	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, nil),
	)
	logger.Infof("user %s logged in from %s", "alice", "192.0.2.1")
	// NOTE: // Output: is intentionally omitted because the template
	// contains {time} and ANSI escape codes, which change per run.
	// The function will still compile and run during `go test`.
}

// ExampleNewLogger_json shows how to configure a JSON output handler with
// custom time formatting and timezone.
func ExampleNewLogger_json() {
	loc, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		panic(err)
	}

	logger := solislog.NewLogger(
		solislog.Extra{
			"service": "api",
			"env":     "dev",
		},
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
			JSON:       true,
			TimeFormat: time.RFC3339,
			Location:   loc,
			// In JSON mode, plain text in the template is ignored.
			// Placeholders define the JSON fields and their order.
			Template: "{time} {level} {message} {extra[service]} {extra[env]} {extra}",
		}),
	)

	logger.Info("json message")
}

// ExampleLogger_Contextualize demonstrates binding extra fields to a logger,
// storing it in a context, and retrieving it in downstream functions.
func ExampleLogger_Contextualize() {
	base := solislog.NewLogger(
		solislog.Extra{"service": "api"},
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "<level>{level}</level> | service={extra[service]} request_id={extra[request_id]} user_id={extra[user_id]} | {message}\n",
		}),
	)

	// Bind request-specific fields
	requestLogger := base.Bind(solislog.Extra{"request_id": "req-123"})

	// Inject into context
	ctx := context.Background()
	ctx = requestLogger.Contextualize(ctx, solislog.Extra{"user_id": "42"})

	// Pass context through layers
	handleRequestCtx(ctx)
}

// Helpers for the context example (unexported, scoped to tests)
func handleRequestCtx(ctx context.Context) {
	logger, ok := solislog.FromContext(ctx)
	if !ok {
		return
	}
	logger.Info("request received")
	processRequestCtx(ctx)
}

func processRequestCtx(ctx context.Context) {
	logger, ok := solislog.FromContext(ctx)
	if !ok {
		return
	}
	logger.Info("processing request")
}

// ExampleHandlerOptions_hooks demonstrates BeforeHook and AfterHook usage,
// per-handler record isolation, and caller metadata injection.
func ExampleHandlerOptions_hooks() {
	logger := solislog.NewLogger(
		solislog.Extra{"service": "api"},
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
			Template:   "<gray>{caller}</gray> | <level>{level}</level> | service={extra[service]} hook={extra[hook]} | {message}\n",
			WithCaller: true,
			BeforeHook: func(record *solislog.Record) {
				record.Message = strings.ToUpper(record.Message)
				record.Extra["hook"] = "before"
			},
			AfterHook: func(record *solislog.Record, msg []byte, successful bool) {
				// Example side effect: count metrics, send rendered output elsewhere,
				// or inspect the final rendered log line.
				_ = msg
			},
		}),
		// Second handler receives the original record (no hooks applied to it)
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, nil),
	)

	logger.Info("hook changed this message only in first handler")
}
