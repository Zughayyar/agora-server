package router

import (
	"agora-server/internal/controller"
	"agora-server/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupRoutes configures all the routes for the API
func SetupRoutes(r *chi.Mux) {
	// Apply middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)

	// Health check route (root route as requested)
	r.Get("/", controller.HealthHandler)

	// API routes with versioning
	r.Route("/api/v1", func(r chi.Router) {
		// API health check
		r.Get("/health", controller.HealthHandler)
	})
}
