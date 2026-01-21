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

func InitLogger() *slog.TextHandler {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	return slog.NewTextHandler(os.Stdout, opts)
}
