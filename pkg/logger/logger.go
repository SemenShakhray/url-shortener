package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Logger struct {
	Log *slog.Logger
}

func SetupLogger(env string) Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return Logger{Log: log}
}

func (l *Logger) Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	slog.Error(msg, args...)
}
