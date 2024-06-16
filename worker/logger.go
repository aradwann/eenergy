package worker

import (
	"context"
	"fmt"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) print(level slog.Level, args ...interface{}) {
	l.logger.LogAttrs(
		context.Background(),
		level,
		fmt.Sprint(args...),
	)
}

// Debug logs a message at Debug level.
func (l *Logger) Debug(args ...interface{}) {
	l.print(slog.LevelDebug, args...)
}

// Info logs a message at Info level.
func (l *Logger) Info(args ...interface{}) {
	l.print(slog.LevelInfo, args...)
}

// Warn logs a message at Warning level.
func (l *Logger) Warn(args ...interface{}) {
	l.print(slog.LevelWarn, args...)
}

// Error logs a message at Error level.
func (l *Logger) Error(args ...interface{}) {
	l.print(slog.LevelError, args...)
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (l *Logger) Fatal(args ...interface{}) {
	l.print(slog.LevelError, args...)
}
