package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, success bool, message string, data any) {
	response := APIResponse{
		Status:  status,
		Success: success,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"status": 500, "success": false, "message":"internal server error"}`, http.StatusInternalServerError)
	}
}
