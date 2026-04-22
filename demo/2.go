package main

import (
	"os"

	"github.com/DasKaroWow/solislog"
)

func example2() {
	logger := solislog.Add(
		os.Stdout,
		solislog.InfoLevel,
		"{time} | {level} | {extra[name]} | {message}\n",
		map[string]string{
			"name": "ivan",
		},
	)

	_ = logger.Info("base logger message")
}
