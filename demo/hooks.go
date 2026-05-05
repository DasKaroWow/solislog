package main

import (
	"os"
	"strings"

	"github.com/DasKaroWow/solislog"
)

func exampleHooksLogger() {
	logger := solislog.NewLogger(
		solislog.Extra{
			"service": "api",
		},
		solislog.NewHandler(os.Stdout, solislog.InfoLevel, &solislog.HandlerOptions{
			Template: "<gray>{caller}</gray> | <level>{level}</level> | service={extra[service]} hook={extra[hook]} | {message}\n",

			BeforeHook: func(record *solislog.Record) {
				record.Message = strings.ToUpper(record.Message)
				record.Extra["hook"] = "before"
			},

			AfterHook: func(record *solislog.Record, rendered string) {
				// Example side effect: count metrics, send rendered output elsewhere,
				// or inspect the final rendered log line.
				_ = rendered
			},
		}),
	)

	logger.Info("hook changed this message")
}
