package router

import (
	"agora-server/internal/handlers"
	"net/http"

	"github.com/uptrace/bun"
)

func SetupRoutes(mux *http.ServeMux, db *bun.DB) {
	// API v1 routes
	apiV1 := http.NewServeMux()
	apiV1.HandleFunc("/health", handlers.HealthHandler)

	// TODO: Add menu item routes when handlers are implemented
	// apiV1.HandleFunc("GET /menu-items", handlers.GetMenuItems(db))
	// apiV1.HandleFunc("POST /menu-items", handlers.CreateMenuItem(db))
	// apiV1.HandleFunc("GET /menu-items/{id}", handlers.GetMenuItem(db))
	// apiV1.HandleFunc("PUT /menu-items/{id}", handlers.UpdateMenuItem(db))
	// apiV1.HandleFunc("DELETE /menu-items/{id}", handlers.DeleteMenuItem(db))

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))

	// Root level health check
	mux.HandleFunc("/health", handlers.HealthHandler)
}
