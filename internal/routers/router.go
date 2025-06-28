package router

import (
	"net/http"

	"github.com/uptrace/bun"

	"github.com/Zughayyar/agora-server/internal/handlers"
)

func SetupRoutes(mux *http.ServeMux, db *bun.DB) {
	// Initialize handlers
	menuItemHandlers := handlers.NewMenuItemHandlers(db)

	// API v1 routes
	apiV1 := http.NewServeMux()

	// Health check routes
	apiV1.HandleFunc("/health", handlers.HealthHandlerWithDB(db))

	// Menu Items CRUD routes
	apiV1.HandleFunc("GET /items", menuItemHandlers.GetAllMenuItems)
	apiV1.HandleFunc("POST /items", menuItemHandlers.CreateMenuItem)
	apiV1.HandleFunc("GET /items/deleted", menuItemHandlers.GetDeletedMenuItems)
	apiV1.HandleFunc("GET /items/category/{category}", menuItemHandlers.GetMenuItemsByCategory)
	apiV1.HandleFunc("GET /items/{id}", menuItemHandlers.GetMenuItemByID)
	apiV1.HandleFunc("PUT /items/{id}", menuItemHandlers.UpdateMenuItem)
	apiV1.HandleFunc("DELETE /items/{id}", menuItemHandlers.DeleteMenuItem)
	apiV1.HandleFunc("POST /items/{id}/restore", menuItemHandlers.RestoreMenuItem)

	// Mount API v1 routes
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))

	// Root level health check (simple, no database dependency)
	mux.HandleFunc("/health", handlers.HealthHandler)
}
