# Validation Package

**Purpose:** Comprehensive input validation, sanitization, and security middleware

## Features

- ✅ SQL Injection prevention
- ✅ XSS (Cross-Site Scripting) prevention  
- ✅ CSRF (Cross-Site Request Forgery) protection
- ✅ Input sanitization
- ✅ Validation middleware
- ✅ Custom validation rules
- ✅ Comprehensive error messages

## Usage

### Basic Validation

```go
import "github.com/delqhi/biometrics/pkg/validation"

v := validation.NewValidator()

// Validate email
err := v.ValidateEmail("user@example.com")

// Validate URL
err := v.ValidateURL("https://example.com")

// Validate UUID
err := v.ValidateUUID("550e8400-e29b-41d4-a716-446655440000")

// Check for SQL injection
if v.ContainsSQLInjection(input) {
    // Reject input
}

// Check for XSS
if v.ContainsXSS(input) {
    // Reject input
}
```

### Sanitization

```go
s := validation.NewSanitizer()

// Sanitize string
clean := s.SanitizeString(input)

// Sanitize HTML
safeHTML := s.SanitizeHTML(input)

// Sanitize URL
safeURL, err := s.SanitizeURL(input)

// Sanitize file path
safePath := s.SanitizeFilePath(input)
```

### Middleware

```go
m := validation.NewMiddleware(nil)

// Chain middleware
handler := validation.ChainMiddleware(
    myHandler,
    m.SecurityHeaders,
    m.CORSMiddleware,
    m.CSRFProtect,
    m.AuthMiddleware,
)

// Add to HTTP server
http.HandleFunc("/api/endpoint", handler)
```

## Configuration

### Custom Validation Rules

```go
v := validation.NewValidator()

v.RegisterCustomRule("custom_rule", func(value interface{}) error {
    // Custom validation logic
    return nil
})
```

### Middleware Config

```go
config := &validation.MiddlewareConfig{
    CSRFEnabled:      true,
    CSRFTokenCookie:  "csrf_token",
    MaxBodySize:      1 << 20, // 1MB
    AllowedOrigins:   []string{"https://example.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
    SkipPaths:        []string{"/health", "/metrics"},
}

m := validation.NewMiddleware(config)
```

## Security Best Practices

1. **Always validate input** - Never trust user input
2. **Use parameterized queries** - Prevent SQL injection at database level
3. **Encode output** - Prevent XSS by encoding HTML output
4. **Use CSRF tokens** - Protect against CSRF attacks
5. **Limit input length** - Prevent buffer overflows
6. **Sanitize before storage** - Clean data before saving
7. **Validate on server** - Client-side validation is not enough

## Testing

```bash
# Run tests
go test ./pkg/validation/... -v

# Run with coverage
go test ./pkg/validation/... -coverprofile=coverage.out

# Run specific test
go test ./pkg/validation/... -run TestValidateEmail -v
```

## Files

- `validation.go` - Core validation logic (400+ lines)
- `sanitization.go` - Input sanitization (300+ lines)
- `middleware.go` - HTTP middleware (250+ lines)
- `validation_test.go` - Unit tests (300+ lines)

## Related Documentation

- [Security Guidelines](../../docs/security/security-guidelines.md)
- [Authentication](../auth/README.md)
- [API Reference](../../docs/api/)
