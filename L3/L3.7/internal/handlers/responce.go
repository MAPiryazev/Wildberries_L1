package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
)

func respondError(w http.ResponseWriter, code int, msg string) {
	resp := models.ErrorResponse{
		Error: msg,
		Code:  code,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
