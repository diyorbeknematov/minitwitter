package traceid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ctxKey string

const traceIDKey ctxKey = "media-service"

func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceIDKey, uuid.NewString())
}

func FromContext(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}
	return "unknown"
}

func TraceServerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		values := md.Get("x-trace-id")

		if len(values) > 0 {
			ctx = context.WithValue(ctx, traceIDKey, values[0])
		}
	}

	return handler(ctx, req)
}
