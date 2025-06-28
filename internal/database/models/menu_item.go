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
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
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
		// Update timestamp on updates
		m.UpdatedAt = time.Now()
	}
	return nil
}

// TableName returns the table name for this model
func (MenuItem) TableName() string {
	return "menu_items"
}

// String returns a string representation of the menu item
func (m *MenuItem) String() string {
	return fmt.Sprintf("MenuItem{ID: %s, Name: %s, Price: %s, Category: %s}",
		m.ID, m.Name, m.Price.String(), m.Category)
}
