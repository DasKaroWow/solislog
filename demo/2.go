package main

import (
	"os"

	"github.com/DasKaroWow/solislog"
)

func example2() {
	logger := solislog.NewLogger(
		solislog.Extra{
			"source": "telegram",
			"id":     "-1", // default value
		},
		solislog.NewHandler(os.Stdout, solislog.InfoLevel,
			&solislog.HandlerOptions{
				Template: "{time} | {level} | {extra[source]} | {extra[id]} | {message}\n",
			},
		),
	)
	logger.Info("logger message1") // source = telegram; id = -1

	boundLogger := logger.Bind(solislog.Extra{"id": "123"})
	boundLogger.Info("logger message2") // source = telegram; id = 123

	logger.Info("logger message3") // source = telegram; id = -1
}
