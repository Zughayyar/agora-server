package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zughayyar/agora-server/internal/database"
	"github.com/Zughayyar/agora-server/internal/middlewares"
	router "github.com/Zughayyar/agora-server/internal/routers"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment variables")
	}

	// Setup structured logger
	var logger *slog.Logger
	if os.Getenv("APP_ENV") == "development" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	slog.SetDefault(logger)

	// Initialize database with connection pooling
	db, err := initDatabase()
	if err != nil {
		logger.Error("Failed to initialize database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			logger.Error("Failed to close database connection", slog.String("error", err.Error()))
		}
	}()

	appName := "Agora Restaurant Management API"
	appVersion := os.Getenv("APP_VERSION")
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "3000" // Updated to match actual usage
	}
	appEnv := os.Getenv("APP_ENV")

	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Setup routes with database dependency
	router.SetupRoutes(mux, db)

	// Add catch-all 404 handler for unmatched routes (except root)
	mux.HandleFunc("/{path...}", middlewares.NotFoundHandler())

	// Apply global middleware stack
	var handler http.Handler = mux
	handler = middlewares.RecoveryMiddleware(handler)
	handler = middlewares.LoggingMiddleware(handler)
	handler = middlewares.CORSMiddleware(handler)

	// Create server with production-ready timeouts
	server := &http.Server{
		Addr:         ":" + appPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine for graceful shutdown
	go func() {
		logger.Info("🚀 Agora Server starting",
			slog.String("app", appName),
			slog.String("version", appVersion),
			slog.String("port", appPort),
			slog.String("env", appEnv),
		)
		logger.Info("🏥 Health endpoints available:",
			slog.String("root", fmt.Sprintf("http://localhost:%s/health", appPort)),
			slog.String("api", fmt.Sprintf("http://localhost:%s/api/v1/health", appPort)),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}

// initDatabase initializes the database connection
func initDatabase() (*bun.DB, error) {
	// Load database configuration from environment
	config := database.LoadConfig()

	// Create database connection with optimized connection pooling
	db, err := database.NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
