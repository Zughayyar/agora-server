package main

import (
	"log/slog"
	"net/http"
	"os"

	"agora-server/internal/router"

	"github.com/go-chi/chi/v5"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	appName := os.Getenv("APP_NAME")
	appVersion := os.Getenv("APP_VERSION")
	appPort := os.Getenv("APP_PORT")
	appEnv := os.Getenv("APP_ENV")

	r := chi.NewRouter()

	router.SetupRoutes(r)

	// Structured logging with context
	logger.Info("Agora Server starting",
		slog.String("app", appName),
		slog.String("version", appVersion),
		slog.String("port", appPort),
		slog.String("env", appEnv),
	)
	logger.Info("API accessible",
		slog.String("url", "http://localhost:"+port),
		slog.String("health_check", "http://localhost:"+port+"/"),
		slog.String("api_health", "http://localhost:"+port+"/api/v1/health"),
	)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Error("Server failed to start",
			slog.String("error", err.Error()),
			slog.String("port", port),
		)
		os.Exit(1)
	}
}
