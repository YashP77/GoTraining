package internal

import (
	"bufio"
	"context"
	"log"
	"log/slog"
	"os"
	"strconv"
)

const key string = "traceID"

func getTraceID(ctx context.Context) string {
	v := ctx.Value(key)
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

func ReadLastTen(ctx context.Context, fileName string) []string {

	file, err := os.Open(fileName)
	if err != nil {
		slog.Error("Error opening file")
		log.Panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) <= 10 {
		slog.Info("File successfully read", "traceID", getTraceID(ctx))
		return lines
	}
	slog.Info("File successfully read", "traceID", getTraceID(ctx))
	return lines[len(lines)-10:]
}
