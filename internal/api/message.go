package api

import (
	"net/http"
)

// MessageResponse is the type returned on error
type MessageResponse struct {
	Messages []string `json:"messages"`
}

// WriteMessage writes errors onto the response writer
func WriteMessage(w http.ResponseWriter, status int, messages ...string) error {
	return WriteJSON(w, status, MessageResponse{
		Messages: messages,
	})
}
