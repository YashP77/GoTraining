package middleware

import (
	"goTraining/internal"
	"net/http"

	"github.com/google/uuid"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Create new traceID with incoming request
		traceID := uuid.NewString()
		ctx := internal.WithTraceID(r.Context(), traceID)

		// Attach ctx to request
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		internal.LogWithTrace(ctx, "TraceID has been attached to request")

	})
}
