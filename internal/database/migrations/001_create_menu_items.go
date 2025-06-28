package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [UP] creating menu_items table...")

		// Create the menu_items table with specified schema
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS menu_items (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				name VARCHAR(100) NOT NULL,
				description TEXT,
				price DECIMAL(10,2) NOT NULL CHECK (price > 0),
				category VARCHAR(50) NOT NULL CHECK (category IN ('appetizer', 'main', 'dessert', 'drink', 'side')),
				is_available BOOLEAN NOT NULL DEFAULT TRUE,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
			);

			-- Create indexes for better query performance
			CREATE INDEX IF NOT EXISTS idx_menu_items_category ON menu_items(category);
			CREATE INDEX IF NOT EXISTS idx_menu_items_is_available ON menu_items(is_available);
			CREATE INDEX IF NOT EXISTS idx_menu_items_created_at ON menu_items(created_at);

			-- Create a trigger to automatically update updated_at
			CREATE OR REPLACE FUNCTION update_updated_at_column()
			RETURNS TRIGGER AS $$
			BEGIN
				NEW.updated_at = CURRENT_TIMESTAMP;
				RETURN NEW;
			END;
			$$ language 'plpgsql';

			DROP TRIGGER IF EXISTS update_menu_items_updated_at ON menu_items;
			CREATE TRIGGER update_menu_items_updated_at
				BEFORE UPDATE ON menu_items
				FOR EACH ROW
				EXECUTE FUNCTION update_updated_at_column();
		`)

		if err != nil {
			return fmt.Errorf("failed to create menu_items table: %w", err)
		}

		fmt.Println(" ✓")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [DOWN] dropping menu_items table...")

		// Drop the table and related objects
		_, err := db.ExecContext(ctx, `
			DROP TRIGGER IF EXISTS update_menu_items_updated_at ON menu_items;
			DROP FUNCTION IF EXISTS update_updated_at_column();
			DROP TABLE IF EXISTS menu_items;
		`)

		if err != nil {
			return fmt.Errorf("failed to drop menu_items table: %w", err)
		}

		fmt.Println(" ✓")
		return nil
	})
}
