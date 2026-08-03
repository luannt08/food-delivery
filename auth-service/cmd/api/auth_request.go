package api

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, ApiError{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

type ApiError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
