package internal

import (
	"context"
	"flag"
)

func CliFlags(ctx context.Context) (string, int) {
	message := flag.String("message", "", "Message from user")
	userID := flag.Int("userID", 0, "ID of user")
	flag.Parse()
	LogWithTrace(ctx, "Command Line Flags set")
	return *message, *userID
}
