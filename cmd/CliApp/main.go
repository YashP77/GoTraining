package main

import (
	"context"
	"goTraining/api"
	"goTraining/internal"
	"goTraining/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	// "github.com/google/uuid"
)

const outputFile = "output/messages.txt"

type key string

const k key = "traceID"

func getTraceID(ctx context.Context) string {
	v := ctx.Value(k)
	if v == nil {
		return ""
	}
	return v.(string)
}

func main() {

	// Create context and add TraceID
	ctx := context.Background()

	log.Printf("Starting server")

	mux := http.NewServeMux()
	mux.Handle("/messages", middleware.TraceMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.CreateMessageHandler(w, *r, outputFile)
	})))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Create signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	internal.ReadLastTen(ctx, outputFile)

	log.Printf("[traceID=%s] %s", getTraceID(ctx), "Application running. Press CTRL+C to exit")
	sig := <-sigChan
	log.Printf("[traceID=%s] %s", getTraceID(ctx), "Received signal: "+sig.String())
	log.Printf("[traceID=%s] %s", getTraceID(ctx), "Shutting down gracefully")

}
