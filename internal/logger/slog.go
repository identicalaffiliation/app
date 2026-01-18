package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger() *Logger {
	h := initLogger()

	return &Logger{
		logger: slog.New(h),
	}
}

func initLogger() *slog.JSONHandler {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	return slog.NewJSONHandler(os.Stdout, opts)
}
