package telemetry

import (
	"context"
	"github.com/google/uuid"
)

type contextKey string

const traceIDKey contextKey = "trace_id"

// InjectTraceID generiert eine neue UUID und legt sie in den Context
func InjectTraceID(ctx context.Context) (context.Context, string) {
	id := uuid.New().String()
	return context.WithValue(ctx, traceIDKey, id), id
}

// GetTraceID holt die ID aus dem Context. Panict NICHT, gibt "unknown" zurück falls leer.
func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}
	return "unknown"
}
