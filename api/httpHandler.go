package api

import (
	"context"
	"encoding/json"
	"goTraining/internal"
	"log/slog"
	"net/http"
)

const OutputFile = "output/messages.txt"
const key string = "traceID"

type createMessageRequest struct {
	Message string `json:"message"`
	UserID  int    `json:"userID"`
}

type createMessageResponse struct {
	TraceID string `json:"traceID"`
	Status  string `json:"status"`
}

func getTraceID(ctx context.Context) string {
	v := ctx.Value(key)
	if v == nil {
		return ""
	}
	return v.(string)
}

// Expects POST /messages with JSON body
func MessageHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Ensure only post methods are allowed
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("invalid request"+err.Error(), "traceID", getTraceID(ctx))
		http.Error(w, "bad request: inavlid JSON", http.StatusBadRequest)
		return
	}

	// Open, write and close file
	file := internal.OpenFile(ctx, OutputFile)
	defer file.Close()
	internal.WriteToFile(ctx, file, req.Message, req.UserID)

	// Perpare and encode response
	resp := createMessageResponse{
		TraceID: getTraceID(ctx),
		Status:  "saved",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
