package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
)

func getTraceID(ctx context.Context) string {
	v := ctx.Value("traceID")
	if v == nil {
		return ""
	}
	return v.(string)
}

func OpenFile(ctx context.Context, fileName string) *os.File {

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("Error creating file")
		log.Panic(err)
	}
	log.Printf("[traceID=%s] %s", getTraceID(ctx), "File successfully created/located")
	return file
}

func WriteToFile(ctx context.Context, file *os.File, message string, userID int) {

	userIDstring := strconv.Itoa(userID)
	_, err := file.WriteString(message + " " + userIDstring + "\n")
	if err != nil {
		slog.Error("Error writing to file")
		log.Panic(err)
	}
	log.Printf("[traceID=%s] %s", getTraceID(ctx), "Message successfully saved in file")
}

func ReadLastTen(ctx context.Context, fileName string) {

	file, err := os.Open(fileName)
	if err != nil {
		slog.Error("Error opening file")
		log.Panic(err)
	}
	defer file.Close()

	log.Printf("[traceID=%s] %s", getTraceID(ctx), "File successfully read")
	fmt.Println("Last 10 messages are:")

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	start := len(lines) - 10
	if start < 0 {
		start = 0
	}

	for _, line := range lines[start:] {
		fmt.Println(line)
	}

	log.Printf("[traceID=%s] %s", getTraceID(ctx), "Last 10 messages are successfully retrieved")
}
