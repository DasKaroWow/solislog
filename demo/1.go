package main

import (
	"os"

	"github.com/DasKaroWow/solislog"
)

func example1() {
	logger := solislog.NewLogger(
		nil,
		solislog.NewHandler(os.Stderr, solislog.InfoLevel, ""), // Handler with default template: "{time} | {level} | {message}\n"
	)

	logger.Info("hello from solislog")
}
