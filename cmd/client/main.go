package main

import (
	"log"
	"log/slog"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	// Build ws:// URL
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/ws-last10",
	}

	log.Println("connecting to", u.String())

	// Dial the websocket server
	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		if resp != nil {
			slog.Error("dial error", "http status", err, "response status", resp.Status)
		}
		slog.Error("dial error", "err", err)
	}
	defer conn.Close()

	log.Println("connected")

	// Read messages until server closes connection
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			// This is expected when the server closes the connection
			slog.Info("connection closed", "err", err)
			return
		}

		switch msgType {
		case websocket.TextMessage:
			slog.Info("received:" + string(msg))
		case websocket.CloseMessage:
			slog.Info("received close frame")
			return
		default:
			slog.Info("received message type", "message type", msgType)
		}
	}
}
