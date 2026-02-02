package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorPayload struct {
	Error string `json:"error"`
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	respondJSON(w, status, errorPayload{Error: msg})
}
