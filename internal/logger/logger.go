package logger

import (
	"io"
	"log/slog"
)

func New(logLevel slog.Level, w io.Writer) {
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
}
