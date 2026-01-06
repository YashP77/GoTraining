package api

import (
	"context"
	"encoding/json"
	"goTraining/internal"
	"log"
	"net/http"
)

type createMessageRequest struct {
	Message string `json:"message"`
	UserID  int    `json:"userID"`
}

type createMessageResponse struct {
	TraceID string `json:"traceID"`
	Status  string `json:"status"`
}

func getTraceID(ctx context.Context) string {
	v := ctx.Value(internal.Key)
	if v == nil {
		return ""
	}
	return v.(string)
}

// Expects POST /messages with JSON body
func CreateMessageHandler(w http.ResponseWriter, r *http.Request, outputFile string) {

	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[traceID=%s] %s", getTraceID(ctx), "invalid request"+err.Error())
		http.Error(w, "bad request: inavlid JSON", http.StatusBadRequest)
		return
	}

	// Open, write and close file
	file := internal.OpenFile(ctx, outputFile)
	defer file.Close()
	internal.WriteToFile(ctx, file, req.Message, req.UserID)
	internal.ReadLastTen(ctx, outputFile)

	// response
	traceID := getTraceID(ctx)
	resp := createMessageResponse{
		TraceID: traceID,
		Status:  "saved",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
