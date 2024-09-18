package util

import (
	"log/slog"
	"os"
	"strings"
)

func NewLogger(lvl, format string) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(lvl) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	default:
		level = slog.LevelError
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	return slog.New(handler)
}
