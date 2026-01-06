package middleware

import (
	"context"
	"goTraining/internal"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func getTraceID(ctx context.Context) string {
	v := ctx.Value(internal.Key)
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
		}

		// Create context and add to request
		ctx := context.WithValue(r.Context(), internal.Key, traceID)
		r = r.WithContext(ctx)
		slog.Info("TraceID has been has been attached to request", "traceID", getTraceID(r.Context()))

		next.ServeHTTP(w, r)

	})
}
