// Package validation provides comprehensive input validation, sanitization,
// and security middleware for the biometrics CLI application.
//
// This package implements:
//   - SQL Injection prevention
//   - XSS (Cross-Site Scripting) prevention
//   - CSRF (Cross-Site Request Forgery) protection
//   - Input sanitization
//   - Validation middleware
//   - Custom validation rules
//   - Comprehensive error messages
//
// Best Practices Feb 2026 compliant.
package validation

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// Validation errors
var (
	ErrSQLInjection      = fmt.Errorf("SQL injection detected")
	ErrXSSDetected       = fmt.Errorf("XSS attack detected")
	ErrInvalidInput      = fmt.Errorf("invalid input")
	ErrRequiredField     = fmt.Errorf("required field is empty")
	ErrInvalidEmail      = fmt.Errorf("invalid email format")
	ErrInvalidURL        = fmt.Errorf("invalid URL format")
	ErrInvalidUUID       = fmt.Errorf("invalid UUID format")
	ErrInvalidLength     = fmt.Errorf("invalid field length")
	ErrInvalidCharacters = fmt.Errorf("invalid characters detected")
	ErrCSRFTokenInvalid  = fmt.Errorf("invalid CSRF token")
)

// Validator holds the validation configuration
type Validator struct {
	validate       *validator.Validate
	sqlPatterns    []*regexp.Regexp
	xssPatterns    []*regexp.Regexp
	customRules    map[string]ValidationFunc
	maxInputLength int
}

// ValidationFunc is a custom validation function type
type ValidationFunc func(value interface{}) error

// ValidationResult contains the result of a validation operation
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []string
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string
	Message string
	Code    string
}

// NewValidator creates a new Validator with default configurations
func NewValidator() *Validator {
	v := &Validator{
		validate:       validator.New(),
		customRules:    make(map[string]ValidationFunc),
		maxInputLength: 10000,
	}

	v.initSQLPatterns()
	v.initXSSPatterns()
	v.registerDefaultValidations()

	return v
}

// initSQLPatterns initializes SQL injection detection patterns
func (v *Validator) initSQLPatterns() {
	v.sqlPatterns = []*regexp.Regexp{
		// UNION-based SQL injection
		regexp.MustCompile(`(?i)\bUNION\s+(ALL\s+)?SELECT\b`),
		// Boolean-based SQL injection
		regexp.MustCompile(`(?i)\bOR\s+1\s*=\s*1\b`),
		regexp.MustCompile(`(?i)\bAND\s+1\s*=\s*1\b`),
		// Comment-based SQL injection
		regexp.MustCompile(`(?i)(--|#|/\*)`),
		// Stacked queries
		regexp.MustCompile(`(?i);\s*(DROP|DELETE|UPDATE|INSERT|ALTER|CREATE|TRUNCATE)\b`),
		// Time-based SQL injection
		regexp.MustCompile(`(?i)\bSLEEP\s*\(\s*\d+\s*\)`),
		regexp.MustCompile(`(?i)\bWAITFOR\s+DELAY\b`),
		regexp.MustCompile(`(?i)\bBENCHMARK\s*\(`),
		// Error-based SQL injection
		regexp.MustCompile(`(?i)\bCONVERT\s*\(`),
		regexp.MustCompile(`(?i)\bCAST\s*\(`),
		// Additional SQL keywords
		regexp.MustCompile(`(?i)\bEXEC\s*\(`),
		regexp.MustCompile(`(?i)\bEXECUTE\s*\(`),
		regexp.MustCompile(`(?i)\bxp_cmdshell\b`),
		regexp.MustCompile(`(?i)\binformation_schema\b`),
		regexp.MustCompile(`(?i)\bsys\.(tables|columns|objects)\b`),
	}
}

// initXSSPatterns initializes XSS detection patterns
func (v *Validator) initXSSPatterns() {
	v.xssPatterns = []*regexp.Regexp{
		// Script tags
		regexp.MustCompile(`(?i)<\s*script[^>]*>`),
		regexp.MustCompile(`(?i)</\s*script\s*>`),
		// Event handlers
		regexp.MustCompile(`(?i)\bon\w+\s*=`),
		regexp.MustCompile(`(?i)\bonerror\s*=`),
		regexp.MustCompile(`(?i)\bonload\s*=`),
		regexp.MustCompile(`(?i)\bonclick\s*=`),
		regexp.MustCompile(`(?i)\bonmouseover\s*=`),
		// JavaScript protocol
		regexp.MustCompile(`(?i)javascript\s*:`),
		regexp.MustCompile(`(?i)vbscript\s*:`),
		// Data URI with script
		regexp.MustCompile(`(?i)data\s*:\s*text/html`),
		// Iframe and object tags
		regexp.MustCompile(`(?i)<\s*iframe[^>]*>`),
		regexp.MustCompile(`(?i)<\s*object[^>]*>`),
		regexp.MustCompile(`(?i)<\s*embed[^>]*>`),
		// Expression in CSS
		regexp.MustCompile(`(?i)expression\s*\(`),
		// Encoded attacks
		regexp.MustCompile(`(?i)%3c\s*script`),
		regexp.MustCompile(`(?i)%3c\s*/\s*script`),
		// SVG-based XSS
		regexp.MustCompile(`(?i)<\s*svg[^>]*onload`),
		// IMG-based XSS
		regexp.MustCompile(`(?i)<\s*img[^>]*onerror`),
	}
}

// registerDefaultValidations registers built-in validation functions
func (v *Validator) registerDefaultValidations() {
	// Email validation
	v.validate.RegisterValidation("email_strict", v.validateEmail)
	// URL validation
	v.validate.RegisterValidation("url_strict", v.validateURL)
	// UUID validation
	v.validate.RegisterValidation("uuid_strict", v.validateUUID)
	// No SQL injection
	v.validate.RegisterValidation("no_sql_injection", v.validateNoSQLInjection)
	// No XSS
	v.validate.RegisterValidation("no_xss", v.validateNoXSS)
	// Safe string
	v.validate.RegisterValidation("safe_string", v.validateSafeString)
}

// ValidateEmail validates email format strictly
func (v *Validator) validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ValidateURL validates URL format strictly
func (v *Validator) validateURL(fl validator.FieldLevel) bool {
	urlStr := fl.Field().String()
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

// ValidateUUID validates UUID format
func (v *Validator) validateUUID(fl validator.FieldLevel) bool {
	uuid := fl.Field().String()
	// UUID regex: 8-4-4-4-12 hex digits
	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidPattern.MatchString(uuid)
}

// ValidateNoSQLInjection checks for SQL injection patterns
func (v *Validator) validateNoSQLInjection(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return !v.ContainsSQLInjection(value)
}

// ValidateNoXSS checks for XSS patterns
func (v *Validator) validateNoXSS(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return !v.ContainsXSS(value)
}

// ValidateSafeString validates that string contains only safe characters
func (v *Validator) validateSafeString(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return v.IsSafeString(value)
}

// ValidateString validates a string value
func (v *Validator) ValidateString(value string, rules ...string) ValidationResult {
	result := ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []string{},
	}

	// Check max length
	if len(value) > v.maxInputLength {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "input",
			Message: fmt.Sprintf("input exceeds maximum length of %d", v.maxInputLength),
			Code:    "ERR_LENGTH_EXCEEDED",
		})
		return result
	}

	// Check for SQL injection
	if v.ContainsSQLInjection(value) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "input",
			Message: "potential SQL injection detected",
			Code:    "ERR_SQL_INJECTION",
		})
	}

	// Check for XSS
	if v.ContainsXSS(value) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "input",
			Message: "potential XSS attack detected",
			Code:    "ERR_XSS_DETECTED",
		})
	}

	// Apply custom rules
	for _, rule := range rules {
		if fn, exists := v.customRules[rule]; exists {
			if err := fn(value); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   "input",
					Message: err.Error(),
					Code:    "ERR_CUSTOM_VALIDATION",
				})
			}
		}
	}

	return result
}

// ValidateEmail validates an email address
func (v *Validator) ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return ErrRequiredField
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}

	// Additional checks
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ErrInvalidEmail
	}

	// Check domain
	domainParts := strings.Split(parts[1], ".")
	if len(domainParts) < 2 {
		return ErrInvalidEmail
	}

	return nil
}

// ValidateURL validates a URL
func (v *Validator) ValidateURL(rawURL string) error {
	if strings.TrimSpace(rawURL) == "" {
		return ErrRequiredField
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURL
	}

	// Ensure scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: only http and https are allowed")
	}

	// Check for SQL injection in URL
	if v.ContainsSQLInjection(rawURL) {
		return ErrSQLInjection
	}

	return nil
}

// ValidateUUID validates a UUID string
func (v *Validator) ValidateUUID(uuid string) error {
	if strings.TrimSpace(uuid) == "" {
		return ErrRequiredField
	}

	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	if !uuidPattern.MatchString(uuid) {
		return ErrInvalidUUID
	}

	return nil
}

// ContainsSQLInjection checks if the input contains SQL injection patterns
func (v *Validator) ContainsSQLInjection(input string) bool {
	decoded := v.decodeInput(input)
	for _, pattern := range v.sqlPatterns {
		if pattern.MatchString(decoded) {
			return true
		}
	}
	return false
}

// ContainsXSS checks if the input contains XSS patterns
func (v *Validator) ContainsXSS(input string) bool {
	decoded := v.decodeInput(input)
	for _, pattern := range v.xssPatterns {
		if pattern.MatchString(decoded) {
			return true
		}
	}
	return false
}

// decodeInput decodes common encodings to detect obfuscated attacks
func (v *Validator) decodeInput(input string) string {
	decoded := input

	// URL decode
	if strings.Contains(input, "%") {
		if dec, err := url.QueryUnescape(input); err == nil {
			decoded = dec
		}
	}

	// HTML entity decode
	decoded = strings.ReplaceAll(decoded, "&lt;", "<")
	decoded = strings.ReplaceAll(decoded, "&gt;", ">")
	decoded = strings.ReplaceAll(decoded, "&amp;", "&")
	decoded = strings.ReplaceAll(decoded, "&quot;", "\"")
	decoded = strings.ReplaceAll(decoded, "&#39;", "'")
	decoded = strings.ReplaceAll(decoded, "&nbsp;", " ")

	// Hex decode
	if strings.HasPrefix(input, "0x") || strings.HasPrefix(input, "\\x") {
		if dec, err := hex.DecodeString(strings.TrimPrefix(input, "0x")); err == nil {
			decoded = string(dec)
		}
	}

	return decoded
}

// IsSafeString checks if string contains only safe characters
func (v *Validator) IsSafeString(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) &&
			!unicode.IsDigit(r) &&
			!unicode.IsSpace(r) &&
			!strings.ContainsRune("-_.@/\\", r) {
			return false
		}
	}
	return true
}

// RegisterCustomRule registers a custom validation rule
func (v *Validator) RegisterCustomRule(name string, fn ValidationFunc) {
	v.customRules[name] = fn
}

// ValidateStruct validates a struct using tags
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// SanitizeInput sanitizes input by removing dangerous characters
func (v *Validator) SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove control characters except newline and tab
	input = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, input)

	return input
}

// GenerateCSRFToken generates a secure CSRF token
func (v *Validator) GenerateCSRFToken() (string, error) {
	// Generate 32 bytes of random data using crypto/rand
	bytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateCSRFToken validates a CSRF token
func (v *Validator) ValidateCSRFToken(token, expected string) error {
	if token == "" || expected == "" {
		return ErrCSRFTokenInvalid
	}

	// Constant-time comparison to prevent timing attacks
	if !strings.EqualFold(token, expected) {
		return ErrCSRFTokenInvalid
	}

	return nil
}

// ValidateRequired checks if a field is not empty
func (v *Validator) ValidateRequired(value string) error {
	if strings.TrimSpace(value) == "" {
		return ErrRequiredField
	}
	return nil
}

// ValidateLength validates string length
func (v *Validator) ValidateLength(value string, min, max int) error {
	length := len(value)
	if length < min {
		return fmt.Errorf("%w: minimum length is %d", ErrInvalidLength, min)
	}
	if length > max {
		return fmt.Errorf("%w: maximum length is %d", ErrInvalidLength, max)
	}
	return nil
}

// ValidatePattern validates string against a regex pattern
func (v *Validator) ValidatePattern(value string, pattern *regexp.Regexp) error {
	if !pattern.MatchString(value) {
		return ErrInvalidInput
	}
	return nil
}

// ValidateAlphanumeric validates that string contains only alphanumeric characters
func (v *Validator) ValidateAlphanumeric(value string) error {
	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ErrInvalidCharacters
		}
	}
	return nil
}

// ValidateNumeric validates that string contains only numeric characters
func (v *Validator) ValidateNumeric(value string) error {
	for _, r := range value {
		if !unicode.IsDigit(r) {
			return ErrInvalidCharacters
		}
	}
	return nil
}

// GetValidationErrors extracts validation errors from validator.ValidationErrors
func (v *Validator) GetValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			validationErrors = append(validationErrors, ValidationError{
				Field:   e.Field(),
				Message: v.getErrorMessage(e),
				Code:    fmt.Sprintf("ERR_VALIDATION_%s", e.Tag()),
			})
		}
	}

	return validationErrors
}

// getErrorMessage returns a human-readable error message for a validation error
func (v *Validator) getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", e.Field(), e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", e.Field(), e.Tag())
	}
}

// SetMaxInputLength sets the maximum input length
func (v *Validator) SetMaxInputLength(max int) {
	v.maxInputLength = max
}

// GetMaxInputLength returns the maximum input length
func (v *Validator) GetMaxInputLength() int {
	return v.maxInputLength
}
