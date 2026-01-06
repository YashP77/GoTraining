package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type key string

const k key = "traceID"

func getTraceID(ctx context.Context) string {
	v := ctx.Value(k)
	if v == nil {
		return ""
	}
	return v.(string)
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get traceID from header or assign new one if empty
		traceID := r.Header.Get("TraceID")
		if traceID == "" {
			traceID = uuid.NewString() + "Generated"
		}

		// Create context and add to request
		ctx := context.WithValue(r.Context(), k, traceID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		log.Printf("[traceID=%s] %s", getTraceID(r.Context()), "TraceID has been attached to request")

	})
}
