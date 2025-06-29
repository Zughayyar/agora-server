package middlewares

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests with response status and timing
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the response writer to capture status code
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     0,
		}

		// Process the request
		next.ServeHTTP(lrw, r)

		// Log the request with all details
		level := slog.LevelInfo

		// Use different log levels based on status code
		switch {
		case lrw.statusCode >= 500:
			level = slog.LevelError
		case lrw.statusCode >= 400:
			level = slog.LevelWarn
		}

		slog.Log(r.Context(), level, "HTTP Request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", lrw.statusCode),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
	})
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// NotFoundHandler returns a professional 404 JSON response
func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SendErrorResponse(w, r, http.StatusNotFound, "Not Found", "Cannot "+r.Method+" "+r.URL.Path)
	}
}

// MethodNotAllowedHandler returns a professional 405 JSON response
func MethodNotAllowedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SendErrorResponse(w, r, http.StatusMethodNotAllowed, "Method Not Allowed", "Method "+r.Method+" is not allowed for "+r.URL.Path)
	}
}

// SendErrorResponse sends a standardized JSON error response
func SendErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, errorType, message string) {
	response := ErrorResponse{
		Message:    message,
		Error:      errorType,
		StatusCode: statusCode,
		Path:       r.URL.Path,
		Timestamp:  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(response); err != nil {
		// Fallback to simple text response if JSON encoding fails
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.Error("Failed to encode error response", slog.String("error", err.Error()))
		return
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		slog.Error("Failed to write error response", slog.String("error", err.Error()))
	}
}

// RecoveryMiddleware recovers from panics and returns a 500 error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered",
					slog.Any("error", err),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
				)
				SendErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", "An unexpected error occurred")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ResponseWriter wrapper to capture status code and response size
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	n, err := lrw.ResponseWriter.Write(b)
	lrw.size += n
	return n, err
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Message    string    `json:"message"`
	Error      string    `json:"error"`
	StatusCode int       `json:"statusCode"`
	Path       string    `json:"path"`
	Timestamp  time.Time `json:"timestamp"`
}
