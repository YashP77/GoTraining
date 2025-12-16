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

	ctx := context.Background()
	ctx = internal.WithTraceID(ctx, uuid.NewString())
	internal.LogWithTrace(ctx, "Starting application")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	message, userID := internal.CliFlags(ctx)
	file := internal.OpenFile(ctx)
	defer file.Close()

	internal.WriteToFile(ctx, *file, message, userID)
	internal.ReadLastTen(ctx)

	internal.LogWithTrace(ctx, "Application running. Press CTRL+C to exit")

	sig := <-sigChan
	internal.LogWithTrace(ctx, "Received signal: "+sig.String())
	internal.LogWithTrace(ctx, "Shutting down gracefully")

}
