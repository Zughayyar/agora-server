package main

import (
	"log/slog"
	"net/http"
	"os"

	"agora-server/internal/middlewares"
	router "agora-server/internal/routers"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Failed to load .env file",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
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

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	appName := "Agora Restaurant Management API"
	appVersion := os.Getenv("APP_VERSION")
	appPort := os.Getenv("APP_PORT")
	appEnv := os.Getenv("APP_ENV")

	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Setup routes
	router.SetupRoutes(mux)

	// Apply global middleware
	var handler http.Handler = mux
	handler = middlewares.LoggingMiddleware(handler)
	handler = middlewares.CORSMiddleware(handler)

	// Structured logging with context
	logger.Info("Agora Server starting",
		slog.String("app", appName),
		slog.String("version", appVersion),
		slog.String("port", appPort),
		slog.String("env", appEnv),
	)

	if err := http.ListenAndServe(":"+appPort, handler); err != nil {
		logger.Error("Server failed to start",
			slog.String("error", err.Error()),
			slog.String("port", appPort),
		)
		os.Exit(1)
	}
}
