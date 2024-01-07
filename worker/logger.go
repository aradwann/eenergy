package worker

import (
	"context"
	"fmt"
	"log/slog"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Print(level slog.Level, args ...interface{}) {
	slog.LogAttrs(
		context.Background(),
		level,
		fmt.Sprint(args...),
	)

}

// Debug logs a message at Debug level.
func (logger *Logger) Debug(args ...interface{}) {
	logger.Print(slog.LevelDebug, args...)
}

// Info logs a message at Info level.
func (logger *Logger) Info(args ...interface{}) {
	logger.Print(slog.LevelInfo, args...)

}

// Warn logs a message at Warning level.
func (logger *Logger) Warn(args ...interface{}) {
	logger.Print(slog.LevelWarn, args...)

}

// Error logs a message at Error level.
func (logger *Logger) Error(args ...interface{}) {
	logger.Print(slog.LevelError, args...)

}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (logger *Logger) Fatal(args ...interface{}) {
	logger.Print(slog.LevelError, args...)

}
