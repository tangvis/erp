package context

import (
	"context"
	"github.com/google/uuid"
)

const (
	TraceIDKey = "__trace_id"
)

func AutoWrapContext(ctx context.Context, traceID string) context.Context {
	// 全局的traceID
	// nolint:staticcheck
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func ForkContext(ctx context.Context) context.Context {
	return AutoWrapContext(context.Background(), GetTraceID(ctx))
}

func GenerateTrace() string {
	return uuid.NewString()
}

func GetTraceID(ctx context.Context) string {
	return ctx.Value(TraceIDKey).(string)
}
