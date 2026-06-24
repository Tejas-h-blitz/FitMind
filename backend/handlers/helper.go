package handlers

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SendJSON(w http.ResponseWriter, status int, success bool, data interface{}, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := APIResponse{
		Success: success,
		Data:    data,
		Error:   errorMsg,
	}
	json.NewEncoder(w).Encode(response)
}
