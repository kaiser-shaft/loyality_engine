package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, msg string, status int) {
	JSON(w, map[string]string{"error": msg}, status)
}
