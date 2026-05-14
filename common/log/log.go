package log

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func Init(service string, level slog.Level) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger = slog.New(handler).With("service", service)
}

func L() *slog.Logger {
	if logger == nil {
		Init("unknown", slog.LevelInfo)
	}
	return logger
}

func Info(msg string, args ...any)  { L().Info(msg, args...) }
func Warn(msg string, args ...any)  { L().Warn(msg, args...) }
func Error(msg string, args ...any) { L().Error(msg, args...) }
func Debug(msg string, args ...any) { L().Debug(msg, args...) }
