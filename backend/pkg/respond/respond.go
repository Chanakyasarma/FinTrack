package respond

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, envelope{"data": data, "success": true})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, envelope{"error": message, "success": false})
}

func ValidationError(w http.ResponseWriter, errors map[string]string) {
	JSON(w, http.StatusUnprocessableEntity, envelope{
		"error":   "validation failed",
		"fields":  errors,
		"success": false,
	})
}
