package internal

import (
	"bufio"
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Helper: create a temp file
func tempFilePath(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("output", "filehandler_test.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	return name
}

func TestOpenFile(t *testing.T) {
	ctx := context.Background()
	path := tempFilePath(t)

	f := OpenFile(ctx, path)
	if f == nil {
		t.Fatalf("OpenFile returned nil")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected file to exist after OpenFile; stat error: %v", err)
	}
	if info.IsDir() {
		t.Fatalf("expected a file, got directory")
	}
	f.Close()
	os.Remove(path)
}

func TestWriteToFile(t *testing.T) {
	ctx := context.Background()
	path := tempFilePath(t)

	f := OpenFile(ctx, path)
	if f == nil {
		t.Fatalf("OpenFile returned nil")
	}

	msg := "hello unit test"
	uid := 42
	WriteToFile(ctx, f, msg, uid)
	f.Close()

	bs, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	got := string(bs)
	expectedSuffix := msg + " " + strconv.Itoa(uid) + "\n"
	if !strings.Contains(got, expectedSuffix) {
		t.Fatalf("expected line not found in file.\nwant suffix: %q\ngot: %q", expectedSuffix, got)
	}
	os.Remove(path)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	// Save original stdout
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Print to stdout
	fn()

	// Restore and read
	w.Close()
	os.Stdout = old

	var sb strings.Builder
	_, err = io.Copy(&sb, r)
	if err != nil {
		t.Fatalf("failed to read pipe: %v", err)
	}
	r.Close()
	return sb.String()
}

func TestReadLastTen(t *testing.T) {
	ctx := context.Background()
	path := tempFilePath(t)

	// Produce 15 lines
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create temp file for writing: %v", err)
	}
	w := bufio.NewWriter(f)
	for i := 1; i <= 15; i++ {
		_, _ = w.WriteString("line " + strconv.Itoa(i) + "\n")
	}
	w.Flush()
	f.Close()

	out := captureStdout(t, func() {
		ReadLastTen(ctx, path)
	})

	if !strings.Contains(out, "Last 10 messages are:") {
		t.Fatalf("expected header in output, got: %q", out)
	}
	if !strings.Contains(out, "line 15") {
		t.Fatalf("expected output to contain last line 'line 15', got: %q", out)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, l := range lines {
		if l == "line 1" {
			t.Fatalf("did not expect 'line 1' to be printed; got: %q", out)
		}
	}

	os.Remove(path)
}
