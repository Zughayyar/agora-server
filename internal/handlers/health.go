package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Zughayyar/agora-server/internal/database"

	"github.com/uptrace/bun"
)

// HealthResponse represents the JSON response for health check
type HealthResponse struct {
	Service   string               `json:"service"`
	Status    string               `json:"status"`
	Timestamp time.Time            `json:"timestamp"`
	Database  DatabaseHealthStatus `json:"database"`
}

// DatabaseHealthStatus represents database health information
type DatabaseHealthStatus struct {
	Status       string        `json:"status"`
	ResponseTime time.Duration `json:"response_time_ms"`
	Error        string        `json:"error,omitempty"`
}

// HealthHandler handles the root route and returns a hello message
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Service:   "agora-server",
		Status:    "healthy",
		Timestamp: time.Now(),
		Database: DatabaseHealthStatus{
			Status: "healthy",
		},
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

// HealthHandlerWithDB handles health check with database connectivity check
func HealthHandlerWithDB(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Service:   "agora-server",
			Status:    "healthy",
			Timestamp: time.Now(),
		}

		// Check database health
		start := time.Now()
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := database.HealthCheck(ctx, db); err != nil {
			response.Database = DatabaseHealthStatus{
				Status:       "unhealthy",
				ResponseTime: time.Since(start),
				Error:        err.Error(),
			}
			response.Status = "degraded" // Overall service is degraded if DB is down
		} else {
			response.Database = DatabaseHealthStatus{
				Status:       "healthy",
				ResponseTime: time.Since(start),
			}
		}

		w.Header().Set("Content-Type", "application/json")

		// Set appropriate HTTP status code
		statusCode := http.StatusOK
		if response.Status == "degraded" {
			statusCode = http.StatusServiceUnavailable
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(statusCode)
		if _, err := w.Write(buf.Bytes()); err != nil {
			slog.Error("Failed to write response body", slog.String("error", err.Error()))
		}
	}
}
