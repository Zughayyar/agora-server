package router

import (
	"net/http"

	"github.com/uptrace/bun"

	"github.com/Zughayyar/agora-server/internal/handlers"
)

// SetupItemRoutes configures all item-related routes
func SetupItemRoutes(mux *http.ServeMux, db *bun.DB) {
	// Initialize handlers
	menuItemHandlers := handlers.NewMenuItemHandlers(db)

	// Menu Items CRUD routes
	mux.HandleFunc("GET /items", menuItemHandlers.GetAllMenuItems)
	mux.HandleFunc("POST /items", menuItemHandlers.CreateMenuItem)
	mux.HandleFunc("GET /items/deleted", menuItemHandlers.GetDeletedMenuItems)
	mux.HandleFunc("GET /items/category/{category}", menuItemHandlers.GetMenuItemsByCategory)
	mux.HandleFunc("GET /items/{id}", menuItemHandlers.GetMenuItemByID)
	mux.HandleFunc("PUT /items/{id}", menuItemHandlers.UpdateMenuItem)
	mux.HandleFunc("DELETE /items/{id}", menuItemHandlers.DeleteMenuItem)
	mux.HandleFunc("POST /items/{id}/restore", menuItemHandlers.RestoreMenuItem)
}
