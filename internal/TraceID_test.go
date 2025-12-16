package internal

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
)

// TestTraceIDEmpty verifies that an empty context returns an empty TraceID.
func TestTraceIDEmpty(t *testing.T) {
	ctx := context.Background()
	if got := TraceID(ctx); got != "" {
		t.Fatalf("expected empty trace id for background context, got %q", got)
	}
}

// TestWithAndTraceID verifies that WithTraceID attaches the ID and TraceID reads it back,
// and that the original context remains unchanged.
func TestWithAndTraceID(t *testing.T) {
	base := context.Background()
	if TraceID(base) != "" {
		t.Fatalf("expected base context to have empty trace id")
	}

	id := "trace-12345"
	ctx := WithTraceID(base, id)
	if got := TraceID(ctx); got != id {
		t.Fatalf("TraceID returned %q; want %q", got, id)
	}

	// Ensure base context still empty (function should not mutate original)
	if TraceID(base) != "" {
		t.Fatalf("expected original base context still to have empty trace id")
	}
}

// TestLogWithTrace captures the standard logger output and ensures LogWithTrace
// includes the trace id and message.
func TestLogWithTrace(t *testing.T) {
	// preserve original output and flags
	origWriter := log.Writer()
	origFlags := log.Flags()
	defer func() {
		log.SetOutput(origWriter)
		log.SetFlags(origFlags)
	}()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0) // simplify matching: remove timestamps etc.

	id := "TID-999"
	msg := "hello trace"
	ctx := WithTraceID(context.Background(), id)

	LogWithTrace(ctx, msg)

	out := buf.String()
	if !strings.Contains(out, "[traceID="+id+"]") {
		t.Fatalf("log output does not contain trace id; got: %q", out)
	}
	if !strings.Contains(out, msg) {
		t.Fatalf("log output does not contain message; got: %q", out)
	}
}
