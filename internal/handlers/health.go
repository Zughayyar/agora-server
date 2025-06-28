package handlers

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// HealthResponse represents the JSON response for health check
type HealthResponse struct {
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthHandler handles the root route and returns a hello message
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Service:   "agora-server",
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(buf.Bytes()); err != nil {
		slog.Error("Failed to write response body", slog.String("error", err.Error()))
	}

}
