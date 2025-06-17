// Package pkg provides utility functions for the application.
package pkg

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// InitLog initializes the logging system with the provided configuration.
func InitLog(level string, format string, file string) error {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // Default log level
	}

	// Set log level with validation
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}

	var w io.Writer
	var err error

	if file == "" {
		w = os.Stdout
	} else {
		w, err = os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //nolint:gosec
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
	}

	var logger *slog.Logger
	// Configure log format
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		h := slog.NewJSONHandler(w, opts)
		logger = slog.New(h)
	case "text":
		h := slog.NewTextHandler(w, opts)
		logger = slog.New(h)
	default:
		h := slog.NewTextHandler(w, opts)
		logger = slog.New(h)
	}

	slog.SetDefault(logger)
	return nil
}
