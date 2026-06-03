package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"game-vault-api/models"
)

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, statusCode int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errorResponse := models.ErrorResponse{
		Code:  code,
		Error: message,
	}
	
	json.NewEncoder(w).Encode(errorResponse)
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// LogError logs an error with context
func LogError(context string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", context, err)
	}
}
