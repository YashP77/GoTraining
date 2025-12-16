package internal

import (
	"bytes"
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Helper to capture stdout
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	// Save original stdout
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Run function that writes to stdout
	fn()

	// Restore and read
	w.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to read stdout pipe: %v", err)
	}
	os.Stdout = old
	return buf.String()
}

func TestOpenFileCreatesFile(t *testing.T) {
	// Ensure output directory exists (OpenFile expects "output/Task2Messages.txt")
	err := os.MkdirAll("output", 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}
	// Remove file if it already exists to test creation
	_ = os.Remove("output/Task2Messages.txt")

	ctx := context.Background()
	f := OpenFile(ctx)
	if f == nil {
		t.Fatalf("OpenFile returned nil file")
	}

	// Ensure file is usable: write a small test string and flush
	_, err = f.WriteString("test\n")
	if err != nil {
		f.Close()
		t.Fatalf("Failed to write to file returned by OpenFile: %v", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}
	// File should now exist
	info, err := os.Stat("output/Task2Messages.txt")
	if err != nil {
		t.Fatalf("Expected file to exist but got error: %v", err)
	}
	if info.IsDir() {
		t.Fatalf("Expected a file but found a directory")
	}
	// Cleanup
	_ = os.Remove("output/Task2Messages.txt")
}

func TestWriteToFileAndReadLastTen_MoreThanTen(t *testing.T) {
	// Prepare environment
	err := os.MkdirAll("output", 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}
	_ = os.Remove("output/Task2Messages.txt")

	ctx := context.Background()
	f := OpenFile(ctx)
	if f == nil {
		t.Fatalf("OpenFile returned nil file")
	}

	// Write 15 lines using WriteToFile
	total := 15
	for i := 1; i <= total; i++ {
		msg := "message" + strconv.Itoa(i)
		WriteToFile(ctx, *f, msg, i)
	}
	// Close file
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close file after writes: %v", err)
	}

	// Capture stdout while calling ReadLastTen
	out := captureStdout(t, func() {
		ReadLastTen(ctx)
	})

	// Split printed lines and trim spaces/newlines
	printed := []string{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if strings.TrimSpace(line) != "" {
			printed = append(printed, strings.TrimSpace(line))
		}
	}

	// Should print last 10 of the 15 lines (i.e., lines 6..15)
	if len(printed) != 10 {
		t.Fatalf("Expected 10 printed lines, got %d: %v", len(printed), printed)
	}

	// Verify that each printed line corresponds to the last 10 writes, the file lines have format: "<message> <userID>"
	for idx, line := range printed {
		expectedUserID := total - 10 + idx + 1 // 6..15
		if !strings.HasSuffix(line, " "+strconv.Itoa(expectedUserID)) {
			t.Fatalf("Printed line %d = %q does not end with expected userID %d", idx, line, expectedUserID)
		}
	}
	// Cleanup
	_ = os.Remove("output/Task2Messages.txt")
}

func TestReadLastTen_FewerThanTen(t *testing.T) {
	// Prepare environment
	err := os.MkdirAll("output", 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}
	_ = os.Remove("output/Task2Messages.txt")

	ctx := context.Background()
	f := OpenFile(ctx)
	if f == nil {
		t.Fatalf("OpenFile returned nil file")
	}

	// Write 5 lines using WriteToFile
	total := 5
	for i := 1; i <= total; i++ {
		msg := "m" + strconv.Itoa(i)
		WriteToFile(ctx, *f, msg, i)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close file after writes: %v", err)
	}

	out := captureStdout(t, func() {
		ReadLastTen(ctx)
	})

	printed := []string{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if strings.TrimSpace(line) != "" {
			printed = append(printed, strings.TrimSpace(line))
		}
	}

	// should print all 5 lines
	if len(printed) != total {
		t.Fatalf("Expected %d printed lines, got %d: %v", total, len(printed), printed)
	}
	// check userIDs 1..5 are present
	for i := 1; i <= total; i++ {
		wantSuffix := " " + strconv.Itoa(i)
		if !strings.HasSuffix(printed[i-1], wantSuffix) {
			t.Fatalf("Line %d = %q does not end with %q", i-1, printed[i-1], wantSuffix)
		}
	}
	// cleanup
	_ = os.Remove("output/Task2Messages.txt")
}
