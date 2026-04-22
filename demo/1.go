package main

import (
	"os"

	"github.com/DasKaroWow/solislog"
)

func example1() {
	logger := solislog.Add(
		os.Stdout,
		solislog.InfoLevel,
		"{time} | {level} | {message}\n",
		nil,
	)

	_ = logger.Info("hello from solislog")
}
