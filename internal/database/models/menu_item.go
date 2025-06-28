package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// MenuItem represents a dish/item on the restaurant menu
type MenuItem struct {
	bun.BaseModel `bun:"table:menu_items,alias:mi"`

	// Primary key - UUID for better distribution and security
	ID uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`

	// Required fields
	Name     string          `bun:"name,notnull" json:"name" validate:"required,min=1,max=100"`
	Price    decimal.Decimal `bun:"price,type:decimal(10,2),notnull" json:"price" validate:"required,gt=0"`
	Category string          `bun:"category,notnull" json:"category" validate:"required,oneof=appetizer main dessert drink side"`

	// Optional fields
	Description *string `bun:"description,type:text" json:"description,omitempty"`
	IsAvailable bool    `bun:"is_available,notnull,default:true" json:"is_available"`

	// Timestamps for auditing
	CreatedAt time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`
}

// BeforeAppendModel is a Bun hook called before inserting/updating
func (m *MenuItem) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		// Set ID if not provided and timestamps
		if m.ID == uuid.Nil {
			m.ID = uuid.New()
		}
		now := time.Now()
		m.CreatedAt = now
		m.UpdatedAt = now
	case *bun.UpdateQuery:
		// Update timestamp on updates (only if not a soft delete)
		if m.DeletedAt == nil {
			m.UpdatedAt = time.Now()
		}
	}
	return nil
}

// SoftDelete marks the record as deleted by setting deleted_at timestamp
func (m *MenuItem) SoftDelete(ctx context.Context, db *bun.DB) error {
	now := time.Now()
	m.DeletedAt = &now
	m.UpdatedAt = now

	_, err := db.NewUpdate().
		Model(m).
		Set("deleted_at = ?", now).
		Set("updated_at = ?", now).
		Where("id = ?", m.ID).
		Exec(ctx)

	return err
}

// Restore restores a soft-deleted record
func (m *MenuItem) Restore(ctx context.Context, db *bun.DB) error {
	m.DeletedAt = nil
	m.UpdatedAt = time.Now()

	_, err := db.NewUpdate().
		Model(m).
		Set("deleted_at = NULL").
		Set("updated_at = ?", m.UpdatedAt).
		Where("id = ?", m.ID).
		Exec(ctx)

	return err
}

// ForceDelete permanently deletes the record from database
func (m *MenuItem) ForceDelete(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().
		Model(m).
		Where("id = ?", m.ID).
		ForceDelete().
		Exec(ctx)

	return err
}

// IsDeleted checks if the record is soft deleted
func (m *MenuItem) IsDeleted() bool {
	return m.DeletedAt != nil
}

// TableName returns the table name for this model
func (MenuItem) TableName() string {
	return "menu_items"
}

// String returns a string representation of the menu item
func (m *MenuItem) String() string {
	status := "active"
	if m.IsDeleted() {
		status = "deleted"
	}
	return fmt.Sprintf("MenuItem{ID: %s, Name: %s, Price: %s, Category: %s, Status: %s}",
		m.ID, m.Name, m.Price.String(), m.Category, status)
}

// MenuItemQuery provides query methods for MenuItem with soft delete support
type MenuItemQuery struct {
	db *bun.DB
}

// NewMenuItemQuery creates a new query builder for MenuItem
func NewMenuItemQuery(db *bun.DB) *MenuItemQuery {
	return &MenuItemQuery{db: db}
}

// All returns all non-deleted menu items
func (q *MenuItemQuery) All(ctx context.Context) ([]MenuItem, error) {
	var items []MenuItem
	err := q.db.NewSelect().
		Model(&items).
		Where("deleted_at IS NULL").
		Scan(ctx)
	return items, err
}

// WithDeleted returns all menu items including soft-deleted ones
func (q *MenuItemQuery) WithDeleted(ctx context.Context) ([]MenuItem, error) {
	var items []MenuItem
	err := q.db.NewSelect().
		Model(&items).
		Scan(ctx)
	return items, err
}

// OnlyDeleted returns only soft-deleted menu items
func (q *MenuItemQuery) OnlyDeleted(ctx context.Context) ([]MenuItem, error) {
	var items []MenuItem
	err := q.db.NewSelect().
		Model(&items).
		Where("deleted_at IS NOT NULL").
		Scan(ctx)
	return items, err
}

// FindByID finds a menu item by ID (excludes soft-deleted)
func (q *MenuItemQuery) FindByID(ctx context.Context, id uuid.UUID) (*MenuItem, error) {
	var item MenuItem
	err := q.db.NewSelect().
		Model(&item).
		Where("id = ? AND deleted_at IS NULL", id).
		Scan(ctx)
	return &item, err
}

// FindByIDWithDeleted finds a menu item by ID (includes soft-deleted)
func (q *MenuItemQuery) FindByIDWithDeleted(ctx context.Context, id uuid.UUID) (*MenuItem, error) {
	var item MenuItem
	err := q.db.NewSelect().
		Model(&item).
		Where("id = ?", id).
		Scan(ctx)
	return &item, err
}
