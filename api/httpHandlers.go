package api

import (
	"encoding/json"
	"goTraining/internal"
	"net/http"
)

type createMessageRequest struct {
	Message string `json:"message"`
	UserID  int    `json:"userID"`
}

type createMessageResponse struct {
	Message string `json:"message"`
}

func CreateMessageHandler(w http.ResponseWriter, r http.Request) {
	ctx := r.Context()

	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
	}
	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	file := internal.OpenFile(ctx, "Task3Messages.txt")
	if err := internal.WriteToFile(ctx, *file, req.Message, req.UserID); err != nil {
		internal.LogWithTrace(ctx, "failed to save message")
		http.Error(w, "failed to save message", http.StatusInternalServerError)
		return
	}

	internal.LogWithTrace(ctx, "message was created via API")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createMessageResponse{Message: "success"})

}
