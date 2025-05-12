package logger

import (
	"log/slog"
	"os"
)

const appName = "cart"
const appVersion = "1.0.0"

var appLogger *slog.Logger

func init() {
	appLogger = slog.New(slog.NewTextHandler(os.Stdout, nil)).
		With("app", appName, "version", appVersion)
}

func Fatal(msg string, args ...any) {
	appLogger.Error(msg, args...)

	os.Exit(1)
}

func Info(msg string, args ...any) {
	appLogger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	appLogger.Error(msg, args...)
}

func Warn(msg string, args ...any) {
	appLogger.Warn(msg, args...)
}

func Debug(msg string, args ...any) {
	appLogger.Debug(msg, args...)
}
