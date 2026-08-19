package traceid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func TraceClientInterceptor(
	ctx context.Context,
	method string,
	req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

	traceID := FromContext(ctx)

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		"x-trace-id",
		traceID,
	)

	return invoker(ctx, method, req, reply, cc, opts...)
}
