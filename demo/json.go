package main

import (
	"os"
	"time"

	"github.com/DasKaroWow/solislog"
)

func exampleJSONLogger() {
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

			// In JSON mode, text is ignored.
			// Placeholders define fields and their order.
			Template: "{time} {level} {message} {extra[service]} {extra[env]} {extra}",
		}),
	)

	logger.Info("json message")
}
