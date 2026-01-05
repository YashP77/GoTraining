package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func getTraceID(ctx context.Context) string {
	v := ctx.Value("traceID")
	if v == nil {
		return ""
	}
	return v.(string)
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		traceID := r.Header.Get("TraceID")
		if traceID == "" {
			traceID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), "traceID", traceID)

		// Attach ctx to request
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		log.Printf("[traceID=%s] %s", getTraceID(r.Context()), "TraceID has been attached to request")

	})
}
