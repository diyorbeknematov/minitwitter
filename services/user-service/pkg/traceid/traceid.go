package traceid

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const traceIDKey ctxKey = "trace_id"

func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceIDKey, uuid.NewString())
}

func FromContext(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}
	return "unknown"
}
