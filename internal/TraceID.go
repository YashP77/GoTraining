package internal

import (
	"context"
	"log"
)

type key string

const k key = "traceID"

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, k, id)
}

func TraceID(ctx context.Context) string {
	v := ctx.Value(k)
	if v == nil {
		return ""
	}
	return v.(string)
}

func LogWithTrace(ctx context.Context, msg string) {
	log.Printf("[traceID=%s] %s", TraceID(ctx), msg)
}
