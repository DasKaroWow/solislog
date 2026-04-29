package main

import (
	"os"
	"time"

	"github.com/DasKaroWow/solislog"
)

func example5() {
	loc, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		panic(err)
	}

	logger := solislog.NewLogger(
		solislog.Extra{
			"source": "telegram",
			"id":     "-1",
		},
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
			Template:   "{time} {level} {message} {extra[id]} {extra}",
			JSON:       true,
			TimeFormat: time.RFC3339,
			Location:   loc,
		}),
	)

	logger.Info("base message")

	requestLogger := logger.Bind(solislog.Extra{
		"id":   "123",
		"path": "/api/users",
	})

	requestLogger.Info("request message")
}
