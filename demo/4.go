package main

import (
	"os"

	"github.com/DasKaroWow/solislog"
)

func example4() {
	logger := solislog.NewLogger(
		solislog.Extra{
			"source": "telegram",
			"id":     "-1",
			"path":   "/unknown",
		},
		solislog.NewHandler(
			os.Stdout,
			solislog.InfoLevel,
			"handler 1 -> {time} | {level} | source={extra[source]} | id={extra[id]} | {message}\n",
		),
		solislog.NewHandler(
			os.Stdout,
			solislog.InfoLevel,
			"handler 2 -> {time} | {level} | source={extra[source]} | path={extra[path]} | {message}\n",
		),
	)

	logger.Info("base message")

	requestLogger := logger.Bind(solislog.Extra{
		"id":   "123",
		"path": "/api/users",
	})

	requestLogger.Info("request message")
}
