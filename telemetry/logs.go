package telemetry

// TODO: config a logger bridge of Otel

// Package logger provides functions to create logger handlers for development and production environments.

import (
	"log/slog"
	"os"

	"github.com/aradwann/eenergy/util"
)

// newDevelopmentLoggerHandler creates a new development logger handler.
// It returns a logger handler configured to output logs in text format to os.Stdout,
// with additional options for development environment.
func newDevelopmentLoggerHandler() slog.Handler {
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource enables adding source code location to log entries.
		AddSource: true,
		// ReplaceAttr is a callback function that replaces specific attributes in log entries.
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// If the attribute is a time key and not associated with any group, remove it.
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	})
}

// newProductionLoggerHandler creates a new production logger handler.
// It returns a logger handler configured to output logs in JSON format to os.Stdout,
// suitable for production environment.
func newProductionLoggerHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource enables adding source code location to log entries.
		AddSource: true,
	})
}

// initLogger initializes the logger based on the environment.
func InitLogger(config util.Config) *slog.Logger {
	var logHandler slog.Handler

	if config.Environment == "development" || config.Environment == "test" {
		logHandler = newDevelopmentLoggerHandler()
	} else {
		logHandler = newProductionLoggerHandler()
	}

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	return logger
}
