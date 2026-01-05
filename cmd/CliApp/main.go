package main

import (
	"context"
	"flag"
	"goTraining/internal"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
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
	ctx = context.WithValue(ctx, k, uuid.NewString())

	// Create signal
	log.Printf("[traceID=%s] %s", getTraceID(ctx), "Starting application")
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
	log.Printf("[traceID=%s] %s", getTraceID(ctx), "Application running. Press CTRL+C to exit")
	sig := <-sigChan
	log.Printf("[traceID=%s] %s", getTraceID(ctx), "Received signal: "+sig.String())
	log.Printf("[traceID=%s] %s", getTraceID(ctx), "Shutting down gracefully")

}
