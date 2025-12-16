package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func errHand(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	message := flag.String("message", "", "Message from user")
	userID := flag.Int("userID", 0, "ID of user")
	flag.Parse()

	userIDstring := strconv.Itoa(*userID)

	messagesFileName := "output/Task1Messages.txt"
	file, err := os.OpenFile(messagesFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errHand(err)
	defer file.Close()

	_, err = file.WriteString(*message + " " + userIDstring + "\n")
	errHand(err)

	fmt.Println("Message successfully saved \n ")
	fmt.Println("Last 10 messages are:")

	f, err := os.Open("output/Task1Messages.txt")
	errHand(err)
	defer f.Close()

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

}
