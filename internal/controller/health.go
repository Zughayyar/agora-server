package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// HealthResponse represents the JSON response for health check
type HealthResponse struct {
	Message   string    `json:"message"`
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthHandler handles the root route and returns a hello message
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Message:   "Hello from Agora Restaurant Management API! 🍽️",
		Service:   "agora-server",
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode health response",
			slog.String("context", "agora-server"),
			slog.String("error", err.Error()),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	slog.Info("Health check completed",
		slog.String("context", "agora-server"),
		slog.String("status", "healthy"),
	)
}
