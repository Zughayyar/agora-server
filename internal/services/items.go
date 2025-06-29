package services

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"

	"github.com/Zughayyar/agora-server/internal/database/models"
)

// MenuItemService handles business logic for menu items
type MenuItemService struct {
	db    *bun.DB
	query *models.MenuItemQuery
}

// NewMenuItemService creates a new menu item service
func NewMenuItemService(db *bun.DB) *MenuItemService {
	return &MenuItemService{
		db:    db,
		query: models.NewMenuItemQuery(db),
	}
}

// CreateMenuItemRequest represents the data needed to create a menu item
type CreateMenuItemRequest struct {
	Name        string          `json:"name" validate:"required,min=1,max=100"`
	Description *string         `json:"description,omitempty"`
	Price       decimal.Decimal `json:"price" validate:"required,gt=0"`
	Category    string          `json:"category" validate:"required,oneof=appetizer main dessert drink side 'fast food'"`
	IsAvailable *bool           `json:"is_available,omitempty"`
}

// UpdateMenuItemRequest represents the data needed to update a menu item
type UpdateMenuItemRequest struct {
	Name        *string          `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string          `json:"description,omitempty"`
	Price       *decimal.Decimal `json:"price,omitempty" validate:"omitempty,gt=0"`
	Category    *string          `json:"category,omitempty" validate:"omitempty,oneof=appetizer main dessert drink side 'fast food'"`
	IsAvailable *bool            `json:"is_available,omitempty"`
}

// MenuItemResponse represents the response structure for menu items
type MenuItemResponse struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Price       decimal.Decimal `json:"price"`
	Category    string          `json:"category"`
	IsAvailable bool            `json:"is_available"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	DeletedAt   *string         `json:"deleted_at,omitempty"`
}

// CreateMenuItem creates a new menu item
func (s *MenuItemService) CreateMenuItem(ctx context.Context, req CreateMenuItemRequest) (*MenuItemResponse, error) {
	// Create new menu item
	item := &models.MenuItem{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		IsAvailable: true, // Default to available
	}

	// Override default if provided
	if req.IsAvailable != nil {
		item.IsAvailable = *req.IsAvailable
	}

	// Insert into database
	_, err := s.db.NewInsert().Model(item).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create menu item: %w", err)
	}

	return s.toResponse(item), nil
}

// GetAllMenuItems retrieves all active (non-deleted) menu items
func (s *MenuItemService) GetAllMenuItems(ctx context.Context) ([]MenuItemResponse, error) {
	items, err := s.query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve menu items: %w", err)
	}

	responses := make([]MenuItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toResponse(&item)
	}

	return responses, nil
}

// GetMenuItemByID retrieves a specific menu item by ID
func (s *MenuItemService) GetMenuItemByID(ctx context.Context, id int) (*MenuItemResponse, error) {
	item, err := s.query.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu item with ID %d: %w", id, err)
	}

	return s.toResponse(item), nil
}

// GetMenuItemsByCategory retrieves menu items by category
func (s *MenuItemService) GetMenuItemsByCategory(ctx context.Context, category string) ([]MenuItemResponse, error) {
	var items []models.MenuItem
	err := s.db.NewSelect().
		Model(&items).
		Where("category = ? AND deleted_at IS NULL", category).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve menu items by category %s: %w", category, err)
	}

	responses := make([]MenuItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toResponse(&item)
	}

	return responses, nil
}

// GetAvailableMenuItems retrieves only available menu items
func (s *MenuItemService) GetAvailableMenuItems(ctx context.Context) ([]MenuItemResponse, error) {
	var items []models.MenuItem
	err := s.db.NewSelect().
		Model(&items).
		Where("is_available = true AND deleted_at IS NULL").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve available menu items: %w", err)
	}

	responses := make([]MenuItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toResponse(&item)
	}

	return responses, nil
}

// UpdateMenuItem updates an existing menu item
func (s *MenuItemService) UpdateMenuItem(ctx context.Context, id int, req UpdateMenuItemRequest) (*MenuItemResponse, error) {
	// First, get the existing item
	item, err := s.query.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu item with ID %d: %w", id, err)
	}

	// Update fields if provided
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = req.Description
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.Category != nil {
		item.Category = *req.Category
	}
	if req.IsAvailable != nil {
		item.IsAvailable = *req.IsAvailable
	}

	// Update in database
	_, err = s.db.NewUpdate().
		Model(item).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to update menu item: %w", err)
	}

	return s.toResponse(item), nil
}

// SoftDeleteMenuItem marks a menu item as deleted (soft delete)
func (s *MenuItemService) SoftDeleteMenuItem(ctx context.Context, id int) error {
	// Get the item first
	item, err := s.query.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find menu item with ID %d: %w", id, err)
	}

	// Perform soft delete
	if err := item.SoftDelete(ctx, s.db); err != nil {
		return fmt.Errorf("failed to soft delete menu item: %w", err)
	}

	return nil
}

// RestoreMenuItem restores a soft-deleted menu item
func (s *MenuItemService) RestoreMenuItem(ctx context.Context, id int) (*MenuItemResponse, error) {
	// Get the item including deleted ones
	item, err := s.query.FindByIDWithDeleted(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find menu item with ID %d: %w", id, err)
	}

	// Check if it's actually deleted
	if !item.IsDeleted() {
		return nil, fmt.Errorf("menu item with ID %d is not deleted", id)
	}

	// Restore the item
	if err := item.Restore(ctx, s.db); err != nil {
		return nil, fmt.Errorf("failed to restore menu item: %w", err)
	}

	return s.toResponse(item), nil
}

// ForceDeleteMenuItem permanently deletes a menu item from database
func (s *MenuItemService) ForceDeleteMenuItem(ctx context.Context, id int) error {
	// Get the item including deleted ones
	item, err := s.query.FindByIDWithDeleted(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find menu item with ID %d: %w", id, err)
	}

	// Permanently delete
	if err := item.ForceDelete(ctx, s.db); err != nil {
		return fmt.Errorf("failed to permanently delete menu item: %w", err)
	}

	return nil
}

// GetDeletedMenuItems retrieves all soft-deleted menu items
func (s *MenuItemService) GetDeletedMenuItems(ctx context.Context) ([]MenuItemResponse, error) {
	items, err := s.query.OnlyDeleted(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve deleted menu items: %w", err)
	}

	responses := make([]MenuItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toResponse(&item)
	}

	return responses, nil
}

// GetAllMenuItemsWithDeleted retrieves all menu items including soft-deleted ones
func (s *MenuItemService) GetAllMenuItemsWithDeleted(ctx context.Context) ([]MenuItemResponse, error) {
	items, err := s.query.WithDeleted(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all menu items: %w", err)
	}

	responses := make([]MenuItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toResponse(&item)
	}

	return responses, nil
}

// SearchMenuItems searches menu items by name or description
func (s *MenuItemService) SearchMenuItems(ctx context.Context, query string) ([]MenuItemResponse, error) {
	var items []models.MenuItem
	searchPattern := "%" + query + "%"

	err := s.db.NewSelect().
		Model(&items).
		Where("(name ILIKE ? OR description ILIKE ?) AND deleted_at IS NULL", searchPattern, searchPattern).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to search menu items: %w", err)
	}

	responses := make([]MenuItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.toResponse(&item)
	}

	return responses, nil
}

// toResponse converts a MenuItem model to MenuItemResponse
func (s *MenuItemService) toResponse(item *models.MenuItem) *MenuItemResponse {
	response := &MenuItemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		Category:    item.Category,
		IsAvailable: item.IsAvailable,
		CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if item.DeletedAt != nil {
		deletedAt := item.DeletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.DeletedAt = &deletedAt
	}

	return response
}
