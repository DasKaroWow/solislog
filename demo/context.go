package main

import (
	"context"
	"os"

	"github.com/DasKaroWow/solislog"
)

func exampleContextLogger() {
	base := solislog.NewLogger(
		solislog.Extra{
			"service": "api",
		},
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "<level>{level}</level> | service={extra[service]} request_id={extra[request_id]} user_id={extra[user_id]} | {message}\n",
		}),
	)

	requestLogger := base.Bind(solislog.Extra{
		"request_id": "req-123",
	})

	ctx := context.Background()
	ctx = requestLogger.Contextualize(ctx, solislog.Extra{
		"user_id": "42",
	})

	handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
	logger, ok := solislog.FromContext(ctx)
	if !ok {
		return
	}

	logger.Info("request received")
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	logger, ok := solislog.FromContext(ctx)
	if !ok {
		return
	}

	logger.Info("processing request")
}
