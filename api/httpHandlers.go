package api

type createMessageRequest struct {
	Message string `json:"message"`
	UserID  int    `json:"userID"`
}
