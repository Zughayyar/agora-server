package router

import (
	"net/http"

	controller "agora-server/internal/handlers"
	middleware "agora-server/internal/middlewares"
)

// SetupRoutes configures all the routes for the API
func SetupRoutes(mux *http.ServeMux) {
	// Create middleware chain
	var handler http.Handler

	// Health check route (root route as requested)
	handler = http.HandlerFunc(controller.HealthHandler)
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	mux.Handle("/", handler)

	// API health check route
	handler = http.HandlerFunc(controller.HealthHandler)
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	mux.Handle("/api/v1/health", handler)
}
