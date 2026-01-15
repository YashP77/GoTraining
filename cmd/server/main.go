package main

import (
	"context"
	"goTraining/api"
	"goTraining/internal"
	"goTraining/middleware"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/gorilla/websocket"
)

const key string = "traceID"

var listTmpl = template.Must(template.New("list").Parse(`
<!doctype html>
<html>
<body>
<h1>Last 10 Messages</h1>
<ul>
{{range .}}
  <li>{{.}}</li>
{{else}}
  <li>No messages</li>
{{end}}
</ul>
</body>
</html>
`))

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {

	internal.StartActor()

	// Create root context
	rootCtx := context.Background()

	// Create channel for shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	slog.Info("Starting server")

	// Create router
	mux := http.NewServeMux()
	mux.Handle("/messages", middleware.TraceMiddleware(http.HandlerFunc(api.MessageHandler)))

	// Static about page
	fs := http.FileServer((http.Dir("cmd/CliApp/static/about")))
	mux.Handle("/about/", http.StripPrefix("/about/", fs))

	// Dynamic list page
	mux.Handle("/list/", middleware.TraceMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lines := internal.ReadLastTen(r.Context(), api.OutputFile)

		w.Header().Set("Content-Type", "text/html")
		_ = listTmpl.Execute(w, lines)
	})))

	// Websocket last 10
	mux.Handle("/ws-last10", middleware.TraceMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Warn("websocket upgrade failed", "err", err)
			return
		}
		defer conn.Close()

		lines := internal.ReadLastTen(r.Context(), api.OutputFile)

		for _, l := range lines {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(l)); err != nil {
				slog.Warn("websocket write failed", "err", err)
				return
			}
		}

		sub := internal.Subscribe()
		defer internal.Unsubscribe(sub)

		// Send incoming published messages until error
		for msg := range sub {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				slog.Warn("websocket write failed", "err", err)
				return
			}
		}
	})))

	// Server config
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return rootCtx
		},
	}

	// Start server asynch
	go func() {
		slog.Info("Server listening on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	slog.Info("Application running. Press CTRL+C to exit")

	// Shutdown
	sig := <-sigChan
	slog.Info("Received signal: " + sig.String())
	slog.Info("Shutting down")

}
