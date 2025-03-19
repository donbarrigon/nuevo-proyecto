package app

import (
	"encoding/json"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"Error","error":"failed to encode response"}`))
	}
}

func ResponseOkJSON(w http.ResponseWriter, data any, status int, message string) {
	jsonMap := map[string]any{
		"message": message,
		"data":    data,
	}
	ResponseJSON(w, jsonMap, http.StatusOK)
}

func ResponseErrorJSON(w http.ResponseWriter, err any, status int, message string) {
	jsonMap := map[string]any{
		"message": message,
		"error":   err,
	}
	ResponseJSON(w, jsonMap, status)
}
