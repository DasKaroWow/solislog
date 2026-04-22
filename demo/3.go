package main

import (
	"context"
	"os"

	"github.com/DasKaroWow/solislog"
)

func example3() {
	base := solislog.Add(
		os.Stdout,
		solislog.InfoLevel,
		"{time} | {level} | {extra[name]} | {extra[id]} | {message}\n",
		map[string]string{
			"name": "ivan",
		},
	)

	ctx := context.Background()
	ctx = base.Contextualize(ctx, map[string]string{
		"id": "0",
	})

	handle(ctx)
}

func handle(ctx context.Context) {
	log, ok := solislog.FromContext(ctx)
	if !ok {
		return
	}

	_ = log.Info("entered handle")
	process(ctx)
}

func process(ctx context.Context) {
	log, ok := solislog.FromContext(ctx)
	if !ok {
		return
	}

	_ = log.Info("processing request")
}
