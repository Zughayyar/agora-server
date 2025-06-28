package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"agora-server/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	// Load database configuration and connect
	config := database.LoadConfig()
	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("🔍 Database Schema Verification")
	fmt.Println("================================")

	// Check if menu_items table exists
	var tableExists bool
	err = db.NewSelect().
		ColumnExpr("EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'menu_items')").
		Scan(ctx, &tableExists)
	if err != nil {
		log.Fatalf("Failed to check table existence: %v", err)
	}

	if tableExists {
		fmt.Println("✅ menu_items table exists")
	} else {
		fmt.Println("❌ menu_items table does not exist")
		return
	}

	// Get table structure
	type ColumnInfo struct {
		ColumnName    string `bun:"column_name"`
		DataType      string `bun:"data_type"`
		IsNullable    string `bun:"is_nullable"`
		ColumnDefault string `bun:"column_default"`
	}

	var columns []ColumnInfo
	err = db.NewRaw(`
		SELECT column_name, data_type, is_nullable, column_default 
		FROM information_schema.columns 
		WHERE table_name = 'menu_items' 
		ORDER BY ordinal_position
	`).Scan(ctx, &columns)

	if err != nil {
		log.Fatalf("Failed to get column info: %v", err)
	}

	fmt.Println("\n📋 Table Structure:")
	fmt.Println("--------------------")
	for _, col := range columns {
		nullable := "NOT NULL"
		if col.IsNullable == "YES" {
			nullable = "NULL"
		}
		defaultVal := ""
		if col.ColumnDefault != "" {
			defaultVal = fmt.Sprintf(" DEFAULT %s", col.ColumnDefault)
		}
		fmt.Printf("%-15s %-20s %-8s%s\n", col.ColumnName, col.DataType, nullable, defaultVal)
	}

	// Check indexes
	type IndexInfo struct {
		IndexName string `bun:"indexname"`
		IndexDef  string `bun:"indexdef"`
	}

	var indexes []IndexInfo
	err = db.NewRaw(`
		SELECT indexname, indexdef 
		FROM pg_indexes 
		WHERE tablename = 'menu_items'
	`).Scan(ctx, &indexes)

	if err != nil {
		log.Fatalf("Failed to get index info: %v", err)
	}

	fmt.Println("\n🔍 Indexes:")
	fmt.Println("------------")
	for _, idx := range indexes {
		fmt.Printf("- %s\n", idx.IndexName)
	}

	// Check triggers
	type TriggerInfo struct {
		TriggerName string `bun:"trigger_name"`
	}

	var triggers []TriggerInfo
	err = db.NewRaw(`
		SELECT trigger_name 
		FROM information_schema.triggers 
		WHERE event_object_table = 'menu_items'
	`).Scan(ctx, &triggers)

	if err != nil {
		log.Fatalf("Failed to get trigger info: %v", err)
	}

	fmt.Println("\n⚡ Triggers:")
	fmt.Println("------------")
	for _, trigger := range triggers {
		fmt.Printf("- %s\n", trigger.TriggerName)
	}

	// Test a simple count
	fmt.Println("\n🧪 Testing Database Access:")
	fmt.Println("----------------------------")

	// Count existing rows
	count, err := db.NewSelect().
		Table("menu_items").
		Count(ctx)
	if err != nil {
		log.Fatalf("Failed to count rows: %v", err)
	}

	fmt.Printf("Current rows in menu_items: %d\n", count)
	fmt.Println("✅ Database verification complete!")
}
