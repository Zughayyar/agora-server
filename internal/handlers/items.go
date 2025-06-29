package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/uptrace/bun"

	"github.com/Zughayyar/agora-server/internal/services"
)

// MenuItemHandlers contains HTTP handlers for menu item operations
type MenuItemHandlers struct {
	service *services.MenuItemService
}

// NewMenuItemHandlers creates a new menu item handlers instance
func NewMenuItemHandlers(db *bun.DB) *MenuItemHandlers {
	return &MenuItemHandlers{
		service: services.NewMenuItemService(db),
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// CreateMenuItem handles POST /api/v1/menu-items
func (h *MenuItemHandlers) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var req services.CreateMenuItemRequest

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Create menu item using service
	item, err := h.service.CreateMenuItem(r.Context(), req)
	if err != nil {
		slog.Error("Failed to create menu item",
			slog.String("error", err.Error()),
			slog.String("name", req.Name),
			slog.String("category", req.Category))
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created item
	h.writeSuccessResponse(w, item, "Menu item created successfully", http.StatusCreated)
}

// GetAllMenuItems handles GET /api/v1/menu-items
func (h *MenuItemHandlers) GetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	// Check query parameters for filtering
	category := r.URL.Query().Get("category")
	availableOnly := r.URL.Query().Get("available") == "true"
	includeDeleted := r.URL.Query().Get("include_deleted") == "true"
	search := r.URL.Query().Get("search")

	var items []services.MenuItemResponse
	var err error

	// Handle different query scenarios
	switch {
	case search != "":
		items, err = h.service.SearchMenuItems(r.Context(), search)
	case category != "":
		items, err = h.service.GetMenuItemsByCategory(r.Context(), category)
	case availableOnly:
		items, err = h.service.GetAvailableMenuItems(r.Context())
	case includeDeleted:
		items, err = h.service.GetAllMenuItemsWithDeleted(r.Context())
	default:
		items, err = h.service.GetAllMenuItems(r.Context())
	}

	if err != nil {
		slog.Error("Failed to retrieve menu items",
			slog.String("error", err.Error()),
			slog.String("category", category),
			slog.Bool("available_only", availableOnly),
			slog.Bool("include_deleted", includeDeleted),
			slog.String("search", search))
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeSuccessResponse(w, items, "Menu items retrieved successfully", http.StatusOK)
}

// GetMenuItemByID handles GET /api/v1/menu-items/{id}
func (h *MenuItemHandlers) GetMenuItemByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.writeErrorResponse(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	// Get menu item by ID
	item, err := h.service.GetMenuItemByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			slog.Warn("Menu item not found", slog.Int("id", id))
			h.writeErrorResponse(w, "Menu item not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to get menu item by ID",
			slog.String("error", err.Error()),
			slog.Int("id", id))
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeSuccessResponse(w, item, "Menu item retrieved successfully", http.StatusOK)
}

// UpdateMenuItem handles PUT /api/v1/menu-items/{id}
func (h *MenuItemHandlers) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.writeErrorResponse(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	// Parse JSON request body
	var req services.UpdateMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Update menu item
	item, err := h.service.UpdateMenuItem(r.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			slog.Warn("Menu item not found for update", slog.Int("id", id))
			h.writeErrorResponse(w, "Menu item not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to update menu item",
			slog.String("error", err.Error()),
			slog.Int("id", id))
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeSuccessResponse(w, item, "Menu item updated successfully", http.StatusOK)
}

// DeleteMenuItem handles DELETE /api/v1/menu-items/{id}
func (h *MenuItemHandlers) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.writeErrorResponse(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	// Check if force delete is requested
	forceDelete := r.URL.Query().Get("force") == "true"

	if forceDelete {
		// Permanently delete
		err = h.service.ForceDeleteMenuItem(r.Context(), id)
	} else {
		// Soft delete
		err = h.service.SoftDeleteMenuItem(r.Context(), id)
	}

	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			slog.Warn("Menu item not found for deletion", slog.Int("id", id))
			h.writeErrorResponse(w, "Menu item not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to delete menu item",
			slog.String("error", err.Error()),
			slog.Int("id", id),
			slog.Bool("force_delete", forceDelete))
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := "Menu item deleted successfully"
	if forceDelete {
		message = "Menu item permanently deleted"
	}

	h.writeSuccessResponse(w, nil, message, http.StatusOK)
}

// RestoreMenuItem handles POST /api/v1/menu-items/{id}/restore
func (h *MenuItemHandlers) RestoreMenuItem(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.writeErrorResponse(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	// Restore menu item
	item, err := h.service.RestoreMenuItem(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			slog.Warn("Menu item not found for restoration", slog.Int("id", id))
			h.writeErrorResponse(w, "Menu item not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "not deleted") {
			slog.Warn("Attempted to restore non-deleted menu item", slog.Int("id", id))
			h.writeErrorResponse(w, "Menu item is not deleted", http.StatusBadRequest)
			return
		}
		slog.Error("Failed to restore menu item",
			slog.String("error", err.Error()),
			slog.Int("id", id))
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeSuccessResponse(w, item, "Menu item restored successfully", http.StatusOK)
}

// GetDeletedMenuItems handles GET /api/v1/menu-items/deleted
func (h *MenuItemHandlers) GetDeletedMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetDeletedMenuItems(r.Context())
	if err != nil {
		slog.Error("Failed to retrieve deleted menu items", slog.String("error", err.Error()))
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeSuccessResponse(w, items, "Deleted menu items retrieved successfully", http.StatusOK)
}

// GetMenuItemsByCategory handles GET /api/v1/items/category/{category}
func (h *MenuItemHandlers) GetMenuItemsByCategory(w http.ResponseWriter, r *http.Request) {
	// Extract category from URL path using Go 1.22+ path value
	category := r.PathValue("category")
	if category == "" {
		h.writeErrorResponse(w, "Category parameter is required", http.StatusBadRequest)
		return
	}

	// Validate category
	validCategories := map[string]bool{
		"appetizer": true,
		"main":      true,
		"dessert":   true,
		"drink":     true,
		"side":      true,
		"fast food": true,
	}

	if !validCategories[category] {
		h.writeErrorResponse(w, "Invalid category. Must be one of: appetizer, main, dessert, drink, side, fast food", http.StatusBadRequest)
		return
	}

	// Get menu items by category
	items, err := h.service.GetMenuItemsByCategory(r.Context(), category)
	if err != nil {
		slog.Error("Failed to retrieve menu items by category",
			slog.String("error", err.Error()),
			slog.String("category", category))
		h.writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeSuccessResponse(w, items, "Menu items retrieved successfully", http.StatusOK)
}

// Helper function to extract ID from URL path
func (h *MenuItemHandlers) extractIDFromPath(path string) (int, error) {
	// Split path and get the last part that should be the ID
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// Find the ID part (should be after "items")
	for i, part := range pathParts {
		if part == "items" && i+1 < len(pathParts) {
			idStr := pathParts[i+1]

			// Skip if this is a special endpoint, not an ID
			if idStr == "restore" || idStr == "deleted" || idStr == "category" {
				return 0, errors.New("invalid ID format: this is a special endpoint")
			}

			// Check if this looks like a restore endpoint: /items/{id}/restore
			if i+2 < len(pathParts) && pathParts[i+2] == "restore" {
				// This is /items/{id}/restore - parse the ID
				return strconv.Atoi(idStr)
			}

			// Check if this is a regular ID endpoint: /items/{id}
			if i+2 >= len(pathParts) {
				// This is just /items/{id}
				return strconv.Atoi(idStr)
			}

			// If we get here, it's an unexpected pattern
			return 0, errors.New("invalid ID format: unexpected path pattern")
		}
	}

	return 0, errors.New("invalid ID format: no ID found in path")
}

// Helper function to write error responses
func (h *MenuItemHandlers) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		// If we can't encode the error response, there's not much we can do
		// The status code has already been set, so the client will get that
		return
	}
}

// Helper function to write success responses
func (h *MenuItemHandlers) writeSuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	successResp := SuccessResponse{
		Data:    data,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(successResp); err != nil {
		// If we can't encode the success response, there's not much we can do
		// The status code has already been set, so the client will get that
		return
	}
}
