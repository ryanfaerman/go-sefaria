package main

import (
	"io"
	"log/slog"
	"os"
	"strings"

	slogpretty "github.com/phsym/console-slog"
)

func NewLogger(format, levelStr, output string) (*slog.Logger, error) {
	var w io.Writer = os.Stderr

	if output != "" {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, err
		}
		w = f
	}
	level := parseLevel(levelStr)

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	case "text":
		handler = slog.NewTextHandler(w, opts)
	case "human", "pretty", "console":
		handler = slogpretty.NewHandler(w, &slogpretty.HandlerOptions{
			Level:      level,
			TimeFormat: "[15:04:05]",
			NoColor:    false,
		})
	default:
		handler = slog.NewTextHandler(w, opts)
	}

	return slog.New(handler), nil
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
