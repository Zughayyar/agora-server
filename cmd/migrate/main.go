package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"agora-server/internal/database"
	"agora-server/internal/database/migrations"

	"github.com/joho/godotenv"
)

func main() {
	// Command line flags
	var (
		action  = flag.String("action", "migrate", "Action to perform: migrate, rollback, status")
		envFile = flag.String("env", ".env", "Environment file to load")
	)
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(*envFile); err != nil {
		slog.Warn(fmt.Sprintf("No %s file found, using system environment variables", *envFile))
	}

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load database configuration
	config := database.LoadConfig()

	// Create database connection
	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Perform the requested action
	switch *action {
	case "migrate", "up":
		slog.Info("Running migrations...")
		if err := migrations.RunMigrations(ctx, db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		slog.Info("✅ Migrations completed successfully")

	case "rollback", "down":
		slog.Info("Rolling back migrations...")
		if err := migrations.RollbackMigrations(ctx, db); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		slog.Info("✅ Rollback completed successfully")

	case "status":
		slog.Info("Checking migration status...")
		if err := migrations.GetMigrationStatus(ctx, db); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions:")
		fmt.Println("  migrate, up    - Run pending migrations")
		fmt.Println("  rollback, down - Rollback last migration")
		fmt.Println("  status         - Show migration status")
		os.Exit(1)
	}
}
