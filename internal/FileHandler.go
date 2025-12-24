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

func OpenFile(ctx context.Context, fileName string) *os.File {

	messagesFileName := "output/" + fileName
	file, err := os.OpenFile(messagesFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errHand(err)
	LogWithTrace(ctx, "File successfully created/located")
	return file

}

func WriteToFile(ctx context.Context, file os.File, message string, userID int) error {

	userIDstring := strconv.Itoa(userID)
	_, err := file.WriteString(message + " " + userIDstring + "\n")
	errHand(err)
	LogWithTrace(ctx, "Message successfully saved in file")
	return err

}

func ReadLastTen(ctx context.Context, fileName string) {

	f, err := os.Open("output/" + fileName)
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

func AppendMessage(ctx context.Context, message string, userID int) error {

	file := OpenFile(ctx, "Task3Messages.txt")
	defer file.Close()

	err := WriteToFile(ctx, *file, message, userID)

	return err

}
