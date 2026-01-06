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

const Key string = "traceID"

func getTraceID(ctx context.Context) string {
	v := ctx.Value(Key)
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
	slog.Info("File successfully created/located", "traceID", getTraceID(ctx))
	return file
}

func WriteToFile(ctx context.Context, file *os.File, message string, userID int) {

	userIDstring := strconv.Itoa(userID)
	_, err := file.WriteString(message + " " + userIDstring + "\n")
	if err != nil {
		slog.Error("Error writing to file")
		log.Panic(err)
	}
	slog.Info("Message successfuly saving in file", "traceID", getTraceID(ctx))
}

func ReadLastTen(ctx context.Context, fileName string) {

	file, err := os.Open(fileName)
	if err != nil {
		slog.Error("Error opening file")
		log.Panic(err)
	}
	defer file.Close()

	slog.Info("File successfully read", "traceID", getTraceID(ctx))
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

	slog.Info("Last 10 messages are successfully retrieved", "traceID", getTraceID(ctx))
}
