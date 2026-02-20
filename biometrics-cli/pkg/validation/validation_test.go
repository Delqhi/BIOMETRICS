package validation

import (
	"regexp"
	"testing"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Fatal("Expected validator, got nil")
	}
	if v.maxInputLength != 10000 {
		t.Errorf("Expected maxInputLength 10000, got %d", v.maxInputLength)
	}
}

func TestValidateEmail(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		email    string
		expected bool
	}{
		{"user@example.com", true},
		{"user.name@example.com", true},
		{"user+tag@example.co.uk", true},
		{"invalid", false},
		{"@example.com", false},
		{"user@", false},
		{"", false},
	}

	for _, tt := range tests {
		err := v.ValidateEmail(tt.email)
		if (err == nil) != tt.expected {
			t.Errorf("ValidateEmail(%q): expected valid=%v, got error=%v", tt.email, tt.expected, err)
		}
	}
}

func TestValidateURL(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		url      string
		expected bool
	}{
		{"https://example.com", true},
		{"http://example.com/path", true},
		{"ftp://example.com", false},
		{"javascript:alert(1)", false},
		{"", false},
	}

	for _, tt := range tests {
		err := v.ValidateURL(tt.url)
		if (err == nil) != tt.expected {
			t.Errorf("ValidateURL(%q): expected valid=%v, got error=%v", tt.url, tt.expected, err)
		}
	}
}

func TestValidateUUID(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		uuid     string
		expected bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550E8400-E29B-41D4-A716-446655440000", true},
		{"invalid-uuid", false},
		{"", false},
	}

	for _, tt := range tests {
		err := v.ValidateUUID(tt.uuid)
		if (err == nil) != tt.expected {
			t.Errorf("ValidateUUID(%q): expected valid=%v, got error=%v", tt.uuid, tt.expected, err)
		}
	}
}

func TestContainsSQLInjection(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		input    string
		expected bool
	}{
		{"SELECT * FROM users", true},
		{"1' OR '1'='1", true},
		{"'; DROP TABLE users; --", true},
		{"UNION SELECT * FROM passwords", true},
		{"normal text", false},
		{"Hello World", false},
	}

	for _, tt := range tests {
		result := v.ContainsSQLInjection(tt.input)
		if result != tt.expected {
			t.Errorf("ContainsSQLInjection(%q): expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestContainsXSS(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		input    string
		expected bool
	}{
		{"<script>alert(1)</script>", true},
		{"<img onerror=alert(1)>", true},
		{"javascript:alert(1)", true},
		{"normal text", false},
		{"<p>Hello</p>", false},
	}

	for _, tt := range tests {
		result := v.ContainsXSS(tt.input)
		if result != tt.expected {
			t.Errorf("ContainsXSS(%q): expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestValidateString(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"normal", "Hello World", true},
		{"sql injection", "SELECT * FROM users", false},
		{"xss", "<script>alert(1)</script>", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.ValidateString(tt.input)
			if result.Valid != tt.valid {
				t.Errorf("Expected valid=%v, got %v", tt.valid, result.Valid)
			}
		})
	}
}

func TestIsSafeString(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		input    string
		expected bool
	}{
		{"HelloWorld123", true},
		{"test@example.com", true},
		{"path/to/file", true},
		{"<script>", false},
		{"'; DROP TABLE", false},
	}

	for _, tt := range tests {
		result := v.IsSafeString(tt.input)
		if result != tt.expected {
			t.Errorf("IsSafeString(%q): expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestValidateLength(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		input    string
		min      int
		max      int
		expected bool
	}{
		{"hello", 1, 10, true},
		{"hi", 5, 10, false},
		{"very long string", 1, 5, false},
	}

	for _, tt := range tests {
		err := v.ValidateLength(tt.input, tt.min, tt.max)
		if (err == nil) != tt.expected {
			t.Errorf("ValidateLength(%q, %d, %d): expected %v, got error=%v", tt.input, tt.min, tt.max, tt.expected, err)
		}
	}
}

func TestValidateAlphanumeric(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		input    string
		expected bool
	}{
		{"Hello123", true},
		{"123456", true},
		{"Hello", true},
		{"Hello!", false},
		{"Hello World", false},
	}

	for _, tt := range tests {
		err := v.ValidateAlphanumeric(tt.input)
		if (err == nil) != tt.expected {
			t.Errorf("ValidateAlphanumeric(%q): expected %v, got error=%v", tt.input, tt.expected, err)
		}
	}
}

func TestGenerateCSRFToken(t *testing.T) {
	v := NewValidator()

	token1, err := v.GenerateCSRFToken()
	if err != nil {
		t.Fatalf("GenerateCSRFToken failed: %v", err)
	}

	token2, err := v.GenerateCSRFToken()
	if err != nil {
		t.Fatalf("GenerateCSRFToken failed: %v", err)
	}

	if token1 == token2 {
		t.Error("Expected different tokens, got same")
	}

	if len(token1) != 64 {
		t.Errorf("Expected token length 64, got %d", len(token1))
	}
}

func TestValidateCSRFToken(t *testing.T) {
	v := NewValidator()

	token := "test-token"

	if err := v.ValidateCSRFToken(token, token); err != nil {
		t.Errorf("Expected valid token, got error: %v", err)
	}

	if err := v.ValidateCSRFToken(token, "different"); err == nil {
		t.Error("Expected error for mismatched tokens")
	}

	if err := v.ValidateCSRFToken("", token); err == nil {
		t.Error("Expected error for empty token")
	}
}

func TestSanitizeInput(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		input    string
		expected string
	}{
		{"hello\x00world", "helloworld"},
		{"  hello  ", "hello"},
		{"hello\nworld", "hello\nworld"},
	}

	for _, tt := range tests {
		result := v.SanitizeInput(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeInput(%q): expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestGetValidationErrors(t *testing.T) {
	v := NewValidator()

	err := v.validate.Var("", "required")
	if err == nil {
		t.Fatal("Expected validation error")
	}

	validationErrors := v.GetValidationErrors(err)
	if len(validationErrors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestSetMaxInputLength(t *testing.T) {
	v := NewValidator()

	v.SetMaxInputLength(5000)
	if v.GetMaxInputLength() != 5000 {
		t.Errorf("Expected max length 5000, got %d", v.GetMaxInputLength())
	}
}

func TestValidatePattern(t *testing.T) {
	v := NewValidator()
	pattern := regexp.MustCompile(`^[a-z]+$`)

	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", true},
		{"Hello", false},
		{"123", false},
	}

	for _, tt := range tests {
		err := v.ValidatePattern(tt.input, pattern)
		if (err == nil) != tt.expected {
			t.Errorf("ValidatePattern(%q): expected %v, got error=%v", tt.input, tt.expected, err)
		}
	}
}

func TestValidateNumeric(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		input    string
		expected bool
	}{
		{"123456", true},
		{"123abc", false},
		{"abc", false},
	}

	for _, tt := range tests {
		err := v.ValidateNumeric(tt.input)
		if (err == nil) != tt.expected {
			t.Errorf("ValidateNumeric(%q): expected %v, got error=%v", tt.input, tt.expected, err)
		}
	}
}

func TestValidateRequired(t *testing.T) {
	v := NewValidator()

	if err := v.ValidateRequired("hello"); err != nil {
		t.Errorf("Expected no error for non-empty string")
	}

	if err := v.ValidateRequired(""); err != ErrRequiredField {
		t.Errorf("Expected ErrRequiredField for empty string")
	}

	if err := v.ValidateRequired("   "); err != ErrRequiredField {
		t.Errorf("Expected ErrRequiredField for whitespace-only string")
	}
}

func TestValidator_DecodeInput(t *testing.T) {
	v := NewValidator()

	tests := []struct {
		input    string
		expected string
	}{
		{"hello%20world", "hello world"},
		{"&lt;script&gt;", "<script>"},
		{"normal", "normal"},
	}

	for _, tt := range tests {
		result := v.decodeInput(tt.input)
		if result != tt.expected {
			t.Errorf("decodeInput(%q): expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestValidator_InitPatterns(t *testing.T) {
	v := NewValidator()

	if len(v.sqlPatterns) == 0 {
		t.Error("Expected SQL patterns to be initialized")
	}

	if len(v.xssPatterns) == 0 {
		t.Error("Expected XSS patterns to be initialized")
	}
}
