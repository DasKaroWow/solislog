package main

import (
	"context"
	"os"

	"github.com/DasKaroWow/solislog"
)

func example3() {
	base := solislog.NewLogger(
		solislog.Extra{
			"source": "telegram",
			"id":     "-1", // default value
		},
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, "{time} | {level} | {extra[source]} | {extra[id]} | {message}\n"),
	)
	base.Info("logger message1") // source = telegram; id = -1

	ctx := context.Background()
	ctx = base.Contextualize(ctx, solislog.Extra{
		"id": "123",
	})
	handle(ctx)

	base.Info("logger message 4")
}

func handle(ctx context.Context) {
	logger, _ := solislog.FromContext(ctx)

	logger.Info("entered handle")
	process(ctx)
}

func process(ctx context.Context) {
	logger, _ := solislog.FromContext(ctx)

	logger.Info("processing request")
}
