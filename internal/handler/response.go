package handler

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]string{"error": message}
	WriteResponse(w, response, statusCode)
}
