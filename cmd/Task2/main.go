package main

import (
	"context"
	"goTraining/internal"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

func main() {

	// Create context and add TraceID
	ctx := context.Background()
	ctx = internal.WithTraceID(ctx, uuid.NewString())

	// Create signal
	internal.LogWithTrace(ctx, "Starting application")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create flags
	message, userID := internal.CliFlags(ctx)
	file := internal.OpenFile(ctx, "Task2Messages.txt")
	defer file.Close()

	// Write and read logic
	internal.WriteToFile(ctx, *file, message, userID)
	internal.ReadLastTen(ctx, "Task2Messages.txt")

	// Shutdown
	internal.LogWithTrace(ctx, "Application running. Press CTRL+C to exit")
	sig := <-sigChan
	internal.LogWithTrace(ctx, "Received signal: "+sig.String())
	internal.LogWithTrace(ctx, "Shutting down gracefully")

}
