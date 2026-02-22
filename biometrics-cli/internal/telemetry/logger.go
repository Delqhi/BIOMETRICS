package telemetry

import (
	"context"
	"log/slog"
	"os"
)

// SetupLogger erzwingt JSON-Logging für maschinelle Auswertung
func SetupLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// LogWithTrace ist ein Helper, der IMMER die TraceID aus dem Context mitloggt
func LogWithTrace(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, attrs ...slog.Attr) {
	traceID := GetTraceID(ctx)
	allAttrs := append([]slog.Attr{slog.String("trace_id", traceID)}, attrs...)
	logger.LogAttrs(ctx, level, msg, allAttrs...)
}
