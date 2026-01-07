package main

import (
	"context"
	"goTraining/api"
	"goTraining/middleware"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const outputFile = "output/messages.txt"
const key string = "traceID"

func getTraceID(ctx context.Context) string {
	v := ctx.Value(key)
	if v == nil {
		return ""
	}
	return v.(string)
}

func main() {

	// Create root context
	rootCtx := context.Background()

	// Create channel for shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	slog.Info("Starting server")

	// Create router
	mux := http.NewServeMux()
	mux.Handle("/messages", middleware.TraceMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.CreateMessageHandler(w, r, outputFile)
	})))

	// Server config
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return rootCtx
		},
	}

	// Start server asynch
	go func() {
		slog.Info("Server listening on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	slog.Info("Application running. Press CTRL+C to exit")

	// Shutdown
	sig := <-sigChan
	slog.Info("Received signal: " + sig.String())
	slog.Info("Shutting down")

}
