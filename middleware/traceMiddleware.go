package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

const key string = "traceID"

func getTraceID(ctx context.Context) string {
	v := ctx.Value(key)
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
			traceID = getTraceID(r.Context())
		}
		if traceID == "" {
			traceID = uuid.NewString()
			slog.Info("TraceID has been has been attached to request", "traceID", getTraceID(r.Context()))
		}

		// Create context and add to request
		ctx := context.WithValue(r.Context(), key, traceID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
