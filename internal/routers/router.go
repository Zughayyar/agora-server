package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/uptrace/bun"

	"github.com/Zughayyar/agora-server/internal/handlers"
)

func SetupRoutes(mux *http.ServeMux, db *bun.DB) {
	// API v1 routes
	apiV1 := http.NewServeMux()

	// Health check routes
	apiV1.HandleFunc("/health", handlers.HealthHandlerWithDB(db))

	// Setup item routes
	SetupItemRoutes(apiV1, db)

	// Mount API v1 routes
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1))

	// Swagger UI - serves at /swagger/
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Root level health check (simple, no database dependency)
	mux.HandleFunc("/health", handlers.HealthHandler)
}
