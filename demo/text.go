package main

import (
	"os"

	"github.com/DasKaroWow/solislog"
)

func exampleTextLogger() {
	logger := solislog.NewLogger(
		solislog.Extra{
			"service": "api",
			"env":     "dev",
		},
		solislog.NewHandler(os.Stdout, solislog.DebugLevel, &solislog.HandlerOptions{
			Template: "<gray>{time}</gray> | <level>{level}</level> | service={extra[service]} env={extra[env]} | {message}\n",
		}),
	)

	logger.Debug("debug message")
	logger.Info("server started")
	logger.Warning("slow request")
	logger.Error("request failed")
}
