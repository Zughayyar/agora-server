package router

import (
	"agora-server/internal/handlers"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/health", handlers.HealthHandler)
}
