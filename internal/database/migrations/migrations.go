package migrations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

// Migrations holds all registered migrations
var Migrations = migrate.NewMigrations()

// RunMigrations runs all pending migrations
func RunMigrations(ctx context.Context, db *bun.DB) error {
	migrator := migrate.NewMigrator(db, Migrations)

	// Initialize migration tables
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	// Run migrations
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if group.IsZero() {
		slog.Info("No new migrations to run")
	} else {
		slog.Info(fmt.Sprintf("Migrated database to %s", group))
	}

	return nil
}

// RollbackMigrations rolls back the last migration group
func RollbackMigrations(ctx context.Context, db *bun.DB) error {
	migrator := migrate.NewMigrator(db, Migrations)

	// Initialize migration tables
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	// Rollback migrations
	group, err := migrator.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	if group.IsZero() {
		slog.Info("No migrations to rollback")
	} else {
		slog.Info(fmt.Sprintf("Rolled back migrations from %s", group))
	}

	return nil
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(ctx context.Context, db *bun.DB) error {
	migrator := migrate.NewMigrator(db, Migrations)

	// Initialize migration tables
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	slog.Info("Migration tables initialized successfully")
	return nil
}
