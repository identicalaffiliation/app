package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	Logger *slog.Logger
}

func NewLogger() *Logger {
	h := InitLogger()

	return &Logger{
		Logger: slog.New(h),
	}
}

func InitLogger() *slog.JSONHandler {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	return slog.NewJSONHandler(os.Stdout, opts)
}
