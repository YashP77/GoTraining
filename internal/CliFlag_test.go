package internal

import (
	"context"
	"flag"
	"io"
	"os"
	"testing"
)

func TestCliFlagsProvided(t *testing.T) {

	// Saves global state
	origCommandLine := flag.CommandLine
	origArgs := os.Args
	// Restore global state
	defer func() {
		flag.CommandLine = origCommandLine
		os.Args = origArgs
	}()

	// Isolated flag parser
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	// Allow test failures
	fs.SetOutput(io.Discard)
	flag.CommandLine = fs

	// Case 1: Flags provided
	os.Args = []string{"cmd", "-message", "TestMessage", "-userID", "1001"}
	ctx := context.Background()
	msg, uid := CliFlags(ctx)

	if msg != "TestMessage" {
		t.Fatalf("Expected message %q, Got %q", "TestMessage", msg)
	}
	if uid != 1001 {
		t.Fatalf("Expected userID %d, Got %d", 1001, uid)
	}

}

func TestCliFlagsEmpty(t *testing.T) {

	// Saves global state
	origCommandLine := flag.CommandLine
	origArgs := os.Args
	// Restore global state
	defer func() {
		flag.CommandLine = origCommandLine
		os.Args = origArgs
	}()

	// Isolated flag parser
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	// Allow test failures
	fs.SetOutput(io.Discard)
	flag.CommandLine = fs

	// Case 2: Flags provided
	os.Args = []string{"cmd"}
	ctx := context.Background()
	msg, uid := CliFlags(ctx)

	if msg != "" {
		t.Fatalf("Expected empty message when not provided, Got %q", msg)
	}
	if uid != 0 {
		t.Fatalf("Expected userID 0 when not provided, Got %d", uid)
	}

}
