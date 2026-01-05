package internal

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
)

func TestTraceID_NoValue(t *testing.T) {
	ctx := context.Background()
	got := TraceID(ctx)
	if got != "" {
		t.Fatalf("expected empty string when no trace ID set, got %q", got)
	}
}

func TestWithTraceID(t *testing.T) {
	origCtx := context.Background()
	if TraceID(origCtx) != "" {
		t.Fatalf("expected original context to have empty TraceID")
	}

	id := "TestTrace-123"
	ctxWith := WithTraceID(origCtx, id)

	// Original context still empty
	if TraceID(origCtx) != "" {
		t.Errorf("original context should remain unchanged; expected empty, got %q", TraceID(origCtx))
	}

	// New context has the id
	if got := TraceID(ctxWith); got != id {
		t.Fatalf("expected TraceID %q, got %q", id, got)
	}
}

func TestLogWithTraceOutput(t *testing.T) {

	// Save and restore previous output
	prev := log.Writer()
	defer log.SetOutput(prev)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	id := "log-TestTrace-1"
	msg := "test"

	ctx := WithTraceID(context.Background(), id)
	LogWithTrace(ctx, msg)

	out := buf.String()

	// Check substrings
	if !strings.Contains(out, "traceID="+id) {
		t.Fatalf("expected log output to contain traceID=%q; got %q", id, out)
	}
	if !strings.Contains(out, msg) {
		t.Fatalf("expected log output to contain message %q; got %q", msg, out)
	}
}
