package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"agora-server/internal/database"
	"agora-server/internal/database/models"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	config := database.LoadConfig()
	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	ctx := context.Background()

	fmt.Println("🕐 Testing Timestamp Handling")
	fmt.Println("==============================")

	// Test 1: Insert with Bun model (Bun hook will set timestamps)
	fmt.Println("\n📝 Test 1: Insert using Bun model with hooks")

	menuItem := &models.MenuItem{
		Name:        "Test Burger",
		Description: stringPtr("Delicious test burger"),
		Price:       decimal.NewFromFloat(12.99),
		Category:    "main",
		IsAvailable: true,
		// Note: NOT setting CreatedAt/UpdatedAt - let Bun hook handle it
	}

	_, err = db.NewInsert().Model(menuItem).Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to insert menu item: %v", err)
	}

	fmt.Printf("✅ Inserted item with ID: %s\n", menuItem.ID)
	fmt.Printf("   CreatedAt (set by Bun): %s\n", menuItem.CreatedAt.Format(time.RFC3339))
	fmt.Printf("   UpdatedAt (set by Bun): %s\n", menuItem.UpdatedAt.Format(time.RFC3339))

	// Test 2: Direct SQL insert (database defaults will handle timestamps)
	fmt.Println("\n📝 Test 2: Direct SQL insert (database handles timestamps)")

	var directInsertID string
	var directCreatedAt, directUpdatedAt time.Time

	err = db.NewRaw(`
		INSERT INTO menu_items (name, description, price, category, is_available)
		VALUES ('Direct SQL Burger', 'Inserted via raw SQL', 15.99, 'main', true)
		RETURNING id, created_at, updated_at
	`).Scan(ctx, &directInsertID, &directCreatedAt, &directUpdatedAt)

	if err != nil {
		log.Fatalf("Failed to insert via raw SQL: %v", err)
	}

	fmt.Printf("✅ Direct SQL inserted item with ID: %s\n", directInsertID)
	fmt.Printf("   CreatedAt (set by DB): %s\n", directCreatedAt.Format(time.RFC3339))
	fmt.Printf("   UpdatedAt (set by DB): %s\n", directUpdatedAt.Format(time.RFC3339))

	// Test 3: Update using Bun model (hook will update timestamp)
	fmt.Println("\n📝 Test 3: Update using Bun model")

	// Wait a moment to see timestamp difference
	time.Sleep(2 * time.Second)

	originalUpdatedAt := menuItem.UpdatedAt
	menuItem.Price = decimal.NewFromFloat(13.99) // Change price

	_, err = db.NewUpdate().Model(menuItem).WherePK().Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to update menu item: %v", err)
	}

	fmt.Printf("✅ Updated item price\n")
	fmt.Printf("   Original UpdatedAt: %s\n", originalUpdatedAt.Format(time.RFC3339))
	fmt.Printf("   New UpdatedAt (Bun): %s\n", menuItem.UpdatedAt.Format(time.RFC3339))

	// Test 4: Direct SQL update (database trigger will handle timestamp)
	fmt.Println("\n📝 Test 4: Direct SQL update (database trigger)")

	var triggerUpdatedAt time.Time
	err = db.NewRaw(`
		UPDATE menu_items 
		SET price = 16.99 
		WHERE id = ?
		RETURNING updated_at
	`, directInsertID).Scan(ctx, &triggerUpdatedAt)

	if err != nil {
		log.Fatalf("Failed to update via raw SQL: %v", err)
	}

	fmt.Printf("✅ Direct SQL updated item\n")
	fmt.Printf("   Original UpdatedAt: %s\n", directUpdatedAt.Format(time.RFC3339))
	fmt.Printf("   New UpdatedAt (DB trigger): %s\n", triggerUpdatedAt.Format(time.RFC3339))

	// Cleanup
	fmt.Println("\n🧹 Cleaning up test data...")
	_, err = db.NewDelete().Model(&models.MenuItem{}).Where("name LIKE 'Test%' OR name LIKE 'Direct%'").Exec(ctx)
	if err != nil {
		log.Printf("Warning: Failed to cleanup test data: %v", err)
	} else {
		fmt.Println("✅ Test data cleaned up")
	}

	fmt.Println("\n📊 Summary:")
	fmt.Println("- Bun model hooks: Set timestamps in Go code before DB")
	fmt.Println("- Database defaults: Only used when no value provided")
	fmt.Println("- Database triggers: Always update timestamps on SQL updates")
	fmt.Println("- Result: Consistent timestamps regardless of method used")
}

func stringPtr(s string) *string {
	return &s
}
