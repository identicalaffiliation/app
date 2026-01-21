package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger() *Logger {
	h := InitUserLogger()

	return &Logger{
		logger: slog.New(h),
	}
}

func InitUserLogger() *slog.JSONHandler {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	return slog.NewJSONHandler(os.Stdout, opts)
}
