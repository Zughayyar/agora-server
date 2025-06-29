package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// Config holds database configuration with connection pool settings
type Config struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string

	// Connection Pool Settings
	MaxOpenConns    int           // Maximum number of open connections
	MaxIdleConns    int           // Maximum number of idle connections
	ConnMaxLifetime time.Duration // Maximum connection lifetime
	ConnMaxIdleTime time.Duration // Maximum connection idle time
}

// LoadConfig loads database configuration from environment variables
func LoadConfig() *Config {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	// Connection pool settings with sensible defaults
	maxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	maxLifetimeMin, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME_MINUTES", "15"))
	maxIdleTimeMin, _ := strconv.Atoi(getEnv("DB_CONN_MAX_IDLE_TIME_MINUTES", "5"))

	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		Database: getEnv("DB_NAME", "agora_db"),
		User:     getEnv("DB_USER", "agora_user"),
		Password: getEnv("DB_PASSWORD", "agora_password"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Connection pool configuration
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: time.Duration(maxLifetimeMin) * time.Minute,
		ConnMaxIdleTime: time.Duration(maxIdleTimeMin) * time.Minute,
	}
}

// NewConnection creates a new Bun database connection with optimized pool settings
func NewConnection(config *Config) (*bun.DB, error) {
	// Build PostgreSQL DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.Database, config.SSLMode,
	)

	// Create underlying SQL connection with pgdriver (Bun's optimized driver)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Configure connection pool for optimal performance
	sqldb.SetMaxOpenConns(config.MaxOpenConns)       // Limit concurrent connections
	sqldb.SetMaxIdleConns(config.MaxIdleConns)       // Keep connections ready
	sqldb.SetConnMaxLifetime(config.ConnMaxLifetime) // Rotate old connections
	sqldb.SetConnMaxIdleTime(config.ConnMaxIdleTime) // Close unused connections

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqldb.PingContext(ctx); err != nil {
		err := sqldb.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create Bun database instance with PostgreSQL dialect
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add debug logging in development mode
	if os.Getenv("APP_ENV") == "development" && os.Getenv("DB_LOG_QUERIES") != "false" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true), // Show full queries
			bundebug.WithEnabled(true), // Enable debugging
		))
	}

	slog.Info("Database connected successfully",
		slog.String("host", config.Host),
		slog.Int("port", config.Port),
		slog.String("database", config.Database),
		slog.Int("max_open_conns", config.MaxOpenConns),
		slog.Int("max_idle_conns", config.MaxIdleConns),
		slog.Duration("conn_max_lifetime", config.ConnMaxLifetime),
	)

	return db, nil
}

// HealthCheck performs a database health check
func HealthCheck(ctx context.Context, db *bun.DB) error {
	// Simple ping with timeout
	return db.PingContext(ctx)
}

// Close gracefully closes the database connection
func Close(db *bun.DB) error {
	return db.Close()
}

// GetStats returns connection pool statistics for monitoring
func GetStats(db *bun.DB) sql.DBStats {
	return db.DB.Stats()
}

// getEnv gets environment variable with fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
