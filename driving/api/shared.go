package api

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error in RealWorld API format
type ErrorResponse struct {
	Errors struct {
		Body []string `json:"body"`
	} `json:"errors"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(messages ...string) ErrorResponse {
	var resp ErrorResponse
	resp.Errors.Body = messages
	return resp
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, NewErrorResponse(message))
}
