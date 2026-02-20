package validation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Middleware provides HTTP middleware for request validation
type Middleware struct {
	validator      *Validator
	sanitizer      *Sanitizer
	csrfEnabled    bool
	csrfCookie     string
	maxBodySize    int64
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
}

// MiddlewareConfig holds configuration for middleware
type MiddlewareConfig struct {
	CSRFEnabled      bool
	CSRFCookieName   string
	MaxBodySize      int64
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	SkipPaths        []string
	CustomValidators map[string]ValidationFunc
}

// DefaultMiddlewareConfig returns default middleware configuration
func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		CSRFEnabled:      true,
		CSRFCookieName:   "csrf_token",
		MaxBodySize:      1 << 20, // 1MB
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		SkipPaths:        []string{"/health", "/metrics"},
		CustomValidators: make(map[string]ValidationFunc),
	}
}

// NewMiddleware creates a new validation middleware
func NewMiddleware(config *MiddlewareConfig) *Middleware {
	if config == nil {
		config = DefaultMiddlewareConfig()
	}

	m := &Middleware{
		validator:      NewValidator(),
		sanitizer:      NewSanitizer(),
		csrfEnabled:    config.CSRFEnabled,
		csrfCookie:     config.CSRFCookieName,
		maxBodySize:    config.MaxBodySize,
		allowedOrigins: config.AllowedOrigins,
		allowedMethods: config.AllowedMethods,
		allowedHeaders: config.AllowedHeaders,
	}

	// Register custom validators
	for name, fn := range config.CustomValidators {
		m.validator.RegisterCustomRule(name, fn)
	}

	return m
}

// ResponseWriter wraps http.ResponseWriter to capture status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// ValidationHandler represents a validated HTTP handler
type ValidationHandler func(http.ResponseWriter, *http.Request, map[string]interface{})

// ValidateRequest validates an HTTP request
func (m *Middleware) ValidateRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip validation for certain paths
		for _, path := range DefaultMiddlewareConfig().SkipPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next(w, r)
				return
			}
		}

		// Validate Content-Type for POST/PUT/PATCH
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") &&
				!strings.Contains(contentType, "application/x-www-form-urlencoded") &&
				!strings.Contains(contentType, "multipart/form-data") {
				http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
				return
			}
		}

		// Limit request body size
		if r.ContentLength > m.maxBodySize {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		next(w, r)
	}
}

// CSRFProtect adds CSRF protection middleware
func (m *Middleware) CSRFProtect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.csrfEnabled {
			next(w, r)
			return
		}

		// Skip CSRF for safe methods
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next(w, r)
			return
		}

		// Validate CSRF token
		token := r.Header.Get("X-CSRF-Token")
		if token == "" {
			token = r.FormValue("csrf_token")
		}

		if token == "" {
			http.Error(w, "CSRF token missing", http.StatusForbidden)
			return
		}

		// Get expected token from cookie
		cookie, err := r.Cookie(m.csrfCookie)
		if err != nil {
			http.Error(w, "CSRF token invalid", http.StatusForbidden)
			return
		}

		// Validate token
		if err := m.validator.ValidateCSRFToken(token, cookie.Value); err != nil {
			http.Error(w, "CSRF token invalid", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// CORSMiddleware adds CORS support
func (m *Middleware) CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		if len(m.allowedOrigins) == 0 {
			// Allow all origins if none specified
			allowed = true
		} else {
			for _, o := range m.allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if origin == "" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.allowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.allowedHeaders, ", "))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// SecurityHeaders adds security headers to response
func (m *Middleware) SecurityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS filter
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")

		// Permissions Policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS (only for HTTPS)
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		next(w, r)
	}
}

// RateLimitMiddleware adds rate limiting (placeholder - use actual rate limiter)
func (m *Middleware) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Placeholder for rate limiting
		// In production, use a proper rate limiter like golang.org/x/time/rate
		next(w, r)
	}
}

// AuthMiddleware adds authentication check (placeholder - integrate with auth package)
func (m *Middleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		// Check Bearer token format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "Token required", http.StatusUnauthorized)
			return
		}

		// Store token in context for downstream handlers
		ctx := context.WithValue(r.Context(), "auth_token", token)
		next(w, r.WithContext(ctx))
	}
}

// ValidateJSONBody validates and sanitizes JSON request body
func (m *Middleware) ValidateJSONBody(schema map[string]string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Body == nil {
				http.Error(w, "Request body required", http.StatusBadRequest)
				return
			}

			// Parse JSON
			var body map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&body); err != nil {
				http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
				return
			}

			// Validate and sanitize each field
			sanitizedBody := make(map[string]interface{})
			for key, value := range body {
				// Check if field is in schema
				rule, exists := schema[key]
				if !exists {
					// Skip unknown fields or return error based on policy
					continue
				}

				// Convert value to string for validation
				strValue, ok := value.(string)
				if !ok {
					sanitizedBody[key] = value
					continue
				}

				// Sanitize
				sanitized := m.sanitizer.SanitizeString(strValue)

				// Validate based on rule
				result := m.validator.ValidateString(sanitized, rule)
				if !result.Valid {
					http.Error(w, fmt.Sprintf("Validation failed for %s: %s", key, result.Errors[0].Message), http.StatusBadRequest)
					return
				}

				sanitizedBody[key] = sanitized
			}

			// Store sanitized body in context
			ctx := context.WithValue(r.Context(), "validated_body", sanitizedBody)
			next(w, r.WithContext(ctx))
		}
	}
}

// GetValidatedBody retrieves validated body from context
func GetValidatedBody(r *http.Request) map[string]interface{} {
	body, ok := r.Context().Value("validated_body").(map[string]interface{})
	if !ok {
		return nil
	}
	return body
}

// GetAuthToken retrieves auth token from context
func GetAuthToken(r *http.Request) string {
	token, ok := r.Context().Value("auth_token").(string)
	if !ok {
		return ""
	}
	return token
}

// GenerateCSRFCookie generates a CSRF token cookie
func (m *Middleware) GenerateCSRFCookie() (*http.Cookie, error) {
	token, err := m.validator.GenerateCSRFToken()
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     m.csrfCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600, // 1 hour
	}

	return cookie, nil
}

// SetCSRFCookie sets CSRF token cookie on response
func (m *Middleware) SetCSRFCookie(w http.ResponseWriter) error {
	cookie, err := m.GenerateCSRFCookie()
	if err != nil {
		return err
	}
	http.SetCookie(w, cookie)
	return nil
}

// ErrorResponse sends a JSON error response
func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   message,
		"code":    fmt.Sprintf("ERR_%d", statusCode),
		"success": "false",
	})
}

// SuccessResponse sends a JSON success response
func SuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    data,
		"success": "true",
	})
}

// ChainMiddleware chains multiple middleware together
func ChainMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// LoggingMiddleware logs request details
func (m *Middleware) LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w)
		next(rw, r)

		// Log request (in production, use proper logging)
		_ = fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, rw.statusCode)
	}
}

// RecoverMiddleware recovers from panics
func (m *Middleware) RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ErrorResponse(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// RequestIDMiddleware adds request ID tracking
func (m *Middleware) RequestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			token, err := m.validator.GenerateCSRFToken()
			if err == nil {
				requestID = token
				w.Header().Set("X-Request-ID", requestID)
			}
		}

		ctx := context.WithValue(r.Context(), "request_id", requestID)
		next(w, r.WithContext(ctx))
	}
}

// GetRequestID retrieves request ID from context
func GetRequestID(r *http.Request) string {
	id, ok := r.Context().Value("request_id").(string)
	if !ok {
		return ""
	}
	return id
}
