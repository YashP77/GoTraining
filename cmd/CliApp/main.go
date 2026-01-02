package main

import (
	"context"
	"flag"
	"goTraining/internal"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

const outputFile = "output/messages.txt"

func main() {

	// Create context and add TraceID
	ctx := context.Background()
	ctx = internal.WithTraceID(ctx, uuid.NewString())

	// Create signal
	internal.LogWithTrace(ctx, "Starting application")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create flags
	message := flag.String("message", "", "Message from user")
	userID := flag.Int("userID", 0, "ID of user")
	flag.Parse()

	// Create or locate file
	file := internal.OpenFile(ctx, outputFile)
	defer file.Close()

	// Write to file
	internal.WriteToFile(ctx, file, *message, *userID)

	// Read and print last 10 lines of file
	internal.ReadLastTen(ctx, outputFile)

	// Shutdown
	internal.LogWithTrace(ctx, "Application running. Press CTRL+C to exit")
	sig := <-sigChan
	internal.LogWithTrace(ctx, "Received signal: "+sig.String())
	internal.LogWithTrace(ctx, "Shutting down gracefully")

}
