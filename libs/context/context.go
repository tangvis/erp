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
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func GenerateTrace() string {
	return uuid.NewString()
}

func GetTranceID(ctx context.Context) string {
	return ctx.Value(TraceIDKey).(string)
}
