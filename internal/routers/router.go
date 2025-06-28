package router

import (
	"agora-server/internal/handlers"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	apiV1 := http.NewServeMux()
	apiV1.HandleFunc("/health", handlers.HealthHandler)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))

	mux.HandleFunc("/health", handlers.HealthHandler)
}
