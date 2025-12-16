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

func errHand(err error) {
	if err != nil {
		slog.Error("Error")
		log.Fatal(err)
	}
}

func OpenFile(ctx context.Context) *os.File {

	messagesFileName := "output/Task2Messages.txt"
	file, err := os.OpenFile(messagesFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errHand(err)
	LogWithTrace(ctx, "File successfully created/located")
	return file

}

func WriteToFile(ctx context.Context, file os.File, message string, userID int) {

	userIDstring := strconv.Itoa(userID)
	_, err := file.WriteString(message + " " + userIDstring + "\n")
	errHand(err)
	LogWithTrace(ctx, "Message successfully saved in file")

}

func ReadLastTen(ctx context.Context) {

	f, err := os.Open("output/Task2Messages.txt")
	errHand(err)
	defer f.Close()

	LogWithTrace(ctx, "File successfully read")

	scanner := bufio.NewScanner(f)
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

	LogWithTrace(ctx, "Last 10 messages are successfully retrieved")

}
