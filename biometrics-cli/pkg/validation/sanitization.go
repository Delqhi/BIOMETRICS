// Package sanitization provides input sanitization utilities for removing
// dangerous content from user input while preserving legitimate data.
//
// This package implements:
//   - HTML entity encoding
//   - Script tag removal
//   - SQL keyword sanitization
//   - File path sanitization
//   - URL sanitization
//   - Unicode normalization
//   - Whitespace normalization
//
// Best Practices Feb 2026 compliant.
package validation

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Sanitizer provides input sanitization functionality
type Sanitizer struct {
	htmlStripPattern     *regexp.Regexp
	scriptStripPattern   *regexp.Regexp
	sqlKeywordPattern    *regexp.Regexp
	pathTraversalPattern *regexp.Regexp
	whitespacePattern    *regexp.Regexp
	controlCharPattern   *regexp.Regexp
	allowedHTMLTags      map[string]bool
	maxRecursionDepth    int
}

// SanitizeConfig holds configuration for sanitization
type SanitizeConfig struct {
	StripHTML           bool
	StripScripts        bool
	StripSQLKeywords    bool
	NormalizePaths      bool
	NormalizeWhitespace bool
	RemoveControlChars  bool
	AllowedHTMLTags     []string
	MaxRecursionDepth   int
}

// DefaultSanitizeConfig returns default sanitization configuration
func DefaultSanitizeConfig() *SanitizeConfig {
	return &SanitizeConfig{
		StripHTML:           true,
		StripScripts:        true,
		StripSQLKeywords:    false,
		NormalizePaths:      true,
		NormalizeWhitespace: true,
		RemoveControlChars:  true,
		AllowedHTMLTags:     []string{},
		MaxRecursionDepth:   3,
	}
}

// NewSanitizer creates a new Sanitizer with default configuration
func NewSanitizer() *Sanitizer {
	return NewSanitizerWithConfig(DefaultSanitizeConfig())
}

// NewSanitizerWithConfig creates a new Sanitizer with custom configuration
func NewSanitizerWithConfig(config *SanitizeConfig) *Sanitizer {
	s := &Sanitizer{
		htmlStripPattern:     regexp.MustCompile(`<[^>]*>`),
		scriptStripPattern:   regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
		sqlKeywordPattern:    regexp.MustCompile(`(?i)\b(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|ALTER|CREATE|TRUNCATE|EXEC|EXECUTE)\b`),
		pathTraversalPattern: regexp.MustCompile(`\.\.[\\/]+`),
		whitespacePattern:    regexp.MustCompile(`\s+`),
		controlCharPattern:   regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`),
		allowedHTMLTags:      make(map[string]bool),
		maxRecursionDepth:    config.MaxRecursionDepth,
	}

	// Populate allowed HTML tags
	for _, tag := range config.AllowedHTMLTags {
		s.allowedHTMLTags[strings.ToLower(tag)] = true
	}

	return s
}

// SanitizeString sanitizes a string input
func (s *Sanitizer) SanitizeString(input string) string {
	return s.SanitizeStringWithConfig(input, DefaultSanitizeConfig())
}

// SanitizeStringWithConfig sanitizes a string with custom configuration
func (s *Sanitizer) SanitizeStringWithConfig(input string, config *SanitizeConfig) string {
	if input == "" {
		return ""
	}

	result := input

	// Remove control characters
	if config.RemoveControlChars {
		result = s.removeControlCharacters(result)
	}

	// Strip scripts
	if config.StripScripts {
		result = s.stripScripts(result)
	}

	// Strip or encode HTML
	if config.StripHTML {
		if len(config.AllowedHTMLTags) > 0 {
			result = s.stripHTMLWithAllowlist(result, config.AllowedHTMLTags)
		} else {
			result = s.stripAllHTML(result)
		}
	}

	// Normalize whitespace
	if config.NormalizeWhitespace {
		result = s.normalizeWhitespace(result)
	}

	// Trim whitespace
	result = strings.TrimSpace(result)

	return result
}

// SanitizeHTML sanitizes HTML input by removing dangerous elements
func (s *Sanitizer) SanitizeHTML(input string) string {
	// First pass: remove scripts
	result := s.stripScripts(input)

	// Second pass: strip dangerous HTML tags
	result = s.stripDangerousHTML(result)

	// Third pass: encode remaining HTML entities
	result = html.EscapeString(result)

	return result
}

// SanitizeURL sanitizes a URL by removing dangerous components
func (s *Sanitizer) SanitizeURL(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	// Only allow http and https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		parsedURL.Scheme = "https"
	}

	// Remove fragments
	parsedURL.Fragment = ""

	// Sanitize query parameters
	if parsedURL.RawQuery != "" {
		query := parsedURL.Query()
		sanitizedQuery := make(url.Values)

		for key, values := range query {
			// Sanitize key
			sanitizedKey := s.SanitizeString(key)
			if sanitizedKey == "" {
				continue
			}

			// Sanitize values
			for _, value := range values {
				sanitizedValue := s.SanitizeString(value)
				if sanitizedValue != "" {
					sanitizedQuery[sanitizedKey] = append(sanitizedQuery[sanitizedKey], sanitizedValue)
				}
			}
		}

		parsedURL.RawQuery = sanitizedQuery.Encode()
	}

	return parsedURL.String(), nil
}

// SanitizeFilePath sanitizes a file path by removing path traversal attempts
func (s *Sanitizer) SanitizeFilePath(inputPath string) string {
	// Remove path traversal patterns
	result := s.pathTraversalPattern.ReplaceAllString(inputPath, "")

	// Clean the path
	result = filepath.Clean(result)

	// Get base name to ensure no directory components
	result = filepath.Base(result)

	return result
}

// SanitizeFileName sanitizes a filename for safe storage
func (s *Sanitizer) SanitizeFileName(filename string) string {
	// Remove path components
	filename = filepath.Base(filename)

	// Remove or replace dangerous characters
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, char, "_")
	}

	// Remove control characters
	filename = s.removeControlCharacters(filename)

	// Limit length
	if len(filename) > 255 {
		ext := filepath.Ext(filename)
		base := filename[:255-len(ext)]
		filename = base + ext
	}

	return filename
}

// SanitizeSQL sanitizes input for safe SQL usage (parameterized queries preferred)
func (s *Sanitizer) SanitizeSQL(input string) string {
	// Escape single quotes
	result := strings.ReplaceAll(input, "'", "''")

	// Remove null bytes
	result = strings.ReplaceAll(result, "\x00", "")

	// Remove dangerous SQL keywords (optional, use parameterized queries instead)
	result = s.sqlKeywordPattern.ReplaceAllString(result, "")

	return result
}

// SanitizeEmail sanitizes an email address
func (s *Sanitizer) SanitizeEmail(email string) string {
	// Trim whitespace
	email = strings.TrimSpace(email)

	// Convert to lowercase
	email = strings.ToLower(email)

	// Remove control characters
	email = s.removeControlCharacters(email)

	// Validate basic format
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}

	// Sanitize local part
	localPart := s.SanitizeString(parts[0])

	// Sanitize domain part
	domainPart := s.SanitizeString(parts[1])

	if localPart == "" || domainPart == "" {
		return ""
	}

	return localPart + "@" + domainPart
}

// SanitizeUsername sanitizes a username
func (s *Sanitizer) SanitizeUsername(username string) string {
	// Trim whitespace
	username = strings.TrimSpace(username)

	// Convert to lowercase
	username = strings.ToLower(username)

	// Replace spaces with underscores
	username = strings.ReplaceAll(username, " ", "_")

	// Remove special characters (allow only alphanumeric and underscore)
	username = regexp.MustCompile(`[^a-z0-9_]`).ReplaceAllString(username, "")

	// Limit length
	if len(username) > 50 {
		username = username[:50]
	}

	return username
}

// SanitizePassword sanitizes a password (minimal sanitization)
func (s *Sanitizer) SanitizePassword(password string) string {
	// Only remove null bytes and control characters
	password = s.removeControlCharacters(password)
	return password
}

// SanitizeJSON sanitizes JSON string input
func (s *Sanitizer) SanitizeJSON(input string) string {
	// Remove control characters except valid JSON whitespace
	result := s.removeControlCharacters(input)

	// Normalize whitespace
	result = s.normalizeWhitespace(result)

	return strings.TrimSpace(result)
}

// SanitizeFilenameForWeb sanitizes a filename for web display
func (s *Sanitizer) SanitizeFilenameForWeb(filename string) string {
	// Sanitize as filename
	filename = s.SanitizeFileName(filename)

	// HTML encode to prevent XSS
	filename = html.EscapeString(filename)

	return filename
}

// SanitizeUserAgent sanitizes a User-Agent string
func (s *Sanitizer) SanitizeUserAgent(ua string) string {
	// Remove control characters
	ua = s.removeControlCharacters(ua)

	// Limit length
	if len(ua) > 500 {
		ua = ua[:500]
	}

	// Strip HTML
	ua = s.stripAllHTML(ua)

	return strings.TrimSpace(ua)
}

// SanitizeLogEntry sanitizes a log entry to prevent log injection
func (s *Sanitizer) SanitizeLogEntry(entry string) string {
	// Remove newlines to prevent log injection
	entry = strings.ReplaceAll(entry, "\n", " ")
	entry = strings.ReplaceAll(entry, "\r", " ")

	// Remove carriage returns
	entry = s.removeControlCharacters(entry)

	// Limit length
	if len(entry) > 10000 {
		entry = entry[:10000]
	}

	return strings.TrimSpace(entry)
}

// SanitizeHTTPHeader sanitizes an HTTP header value
func (s *Sanitizer) SanitizeHTTPHeader(value string) string {
	// Remove newlines to prevent header injection
	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\r", "")

	// Remove control characters
	value = s.removeControlCharacters(value)

	// Limit length
	if len(value) > 1000 {
		value = value[:1000]
	}

	return strings.TrimSpace(value)
}

// stripScripts removes script tags and their content
func (s *Sanitizer) stripScripts(input string) string {
	return s.scriptStripPattern.ReplaceAllString(input, "")
}

// stripAllHTML removes all HTML tags
func (s *Sanitizer) stripAllHTML(input string) string {
	return s.htmlStripPattern.ReplaceAllString(input, "")
}

// stripDangerousHTML removes dangerous HTML tags while keeping safe ones
func (s *Sanitizer) stripDangerousHTML(input string) string {
	dangerousTags := []string{
		"script", "iframe", "object", "embed", "frame", "frameset",
		"applet", "layer", "link", "meta", "style", "base",
	}

	result := input
	for _, tag := range dangerousTags {
		pattern := regexp.MustCompile(`(?i)<` + tag + `[^>]*>.*?</` + tag + `>`)
		result = pattern.ReplaceAllString(result, "")

		// Also remove self-closing tags
		pattern = regexp.MustCompile(`(?i)<` + tag + `[^>]*/?>`)
		result = pattern.ReplaceAllString(result, "")
	}

	return result
}

// stripHTMLWithAllowlist removes HTML tags except those in the allowlist
func (s *Sanitizer) stripHTMLWithAllowlist(input string, allowedTags []string) string {
	// Build allowlist map
	allowlist := make(map[string]bool)
	for _, tag := range allowedTags {
		allowlist[strings.ToLower(tag)] = true
	}

	// Find all tags
	tagPattern := regexp.MustCompile(`<(/?)(\w+)[^>]*>`)
	result := tagPattern.ReplaceAllStringFunc(input, func(match string) string {
		matches := tagPattern.FindStringSubmatch(match)
		if len(matches) < 3 {
			return match
		}

		tagName := strings.ToLower(matches[2])
		if allowlist[tagName] {
			return match
		}
		return ""
	})

	return result
}

// removeControlCharacters removes control characters from string
func (s *Sanitizer) removeControlCharacters(input string) string {
	return s.controlCharPattern.ReplaceAllString(input, "")
}

// normalizeWhitespace normalizes whitespace characters
func (s *Sanitizer) normalizeWhitespace(input string) string {
	return s.whitespacePattern.ReplaceAllString(input, " ")
}

// NormalizeUnicode normalizes Unicode string to NFC form
func (s *Sanitizer) NormalizeUnicode(input string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	reader := transform.NewReader(strings.NewReader(input), t)
	result, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// SanitizeAndValidate performs sanitization followed by validation
func (s *Sanitizer) SanitizeAndValidate(input string, v *Validator) (string, error) {
	// Sanitize
	sanitized := s.SanitizeString(input)

	// Validate
	result := v.ValidateString(sanitized)
	if !result.Valid && len(result.Errors) > 0 {
		return "", fmt.Errorf("validation failed: %s", result.Errors[0].Message)
	}

	return sanitized, nil
}

// BatchSanitize sanitizes multiple strings
func (s *Sanitizer) BatchSanitize(inputs []string) []string {
	results := make([]string, len(inputs))
	for i, input := range inputs {
		results[i] = s.SanitizeString(input)
	}
	return results
}

// SanitizeMap sanitizes all string values in a map
func (s *Sanitizer) SanitizeMap(input map[string]string) map[string]string {
	result := make(map[string]string)
	for key, value := range input {
		sanitizedKey := s.SanitizeString(key)
		sanitizedValue := s.SanitizeString(value)
		if sanitizedKey != "" {
			result[sanitizedKey] = sanitizedValue
		}
	}
	return result
}

// SanitizeSlice sanitizes all strings in a slice
func (s *Sanitizer) SanitizeSlice(input []string) []string {
	result := make([]string, len(input))
	for i, value := range input {
		result[i] = s.SanitizeString(value)
	}
	return result
}

// EscapeHTML escapes HTML special characters
func (s *Sanitizer) EscapeHTML(input string) string {
	return html.EscapeString(input)
}

// UnescapeHTML unescapes HTML special characters
func (s *Sanitizer) UnescapeHTML(input string) string {
	return html.UnescapeString(input)
}

// SanitizeRichText sanitizes rich text while preserving formatting
func (s *Sanitizer) SanitizeRichText(input string) string {
	// Preserve common formatting tags
	allowedTags := []string{"p", "br", "strong", "em", "u", "ul", "ol", "li", "h1", "h2", "h3", "h4", "h5", "h6", "blockquote", "code", "pre"}

	result := s.stripHTMLWithAllowlist(input, allowedTags)
	result = s.stripScripts(result)
	result = s.normalizeWhitespace(result)

	return strings.TrimSpace(result)
}

// SanitizeMarkdown sanitizes Markdown content
func (s *Sanitizer) SanitizeMarkdown(input string) string {
	// Remove HTML from Markdown
	input = s.stripAllHTML(input)

	// Remove JavaScript links
	input = regexp.MustCompile(`(?i)\[([^\]]+)\]\s*\(\s*javascript:[^)]*\)`).ReplaceAllString(input, "")

	// Remove data: URLs in images
	input = regexp.MustCompile(`(?i)!\[([^\]]*)\]\s*\(\s*data:[^)]*\)`).ReplaceAllString(input, "")

	return input
}

// SanitizeCSVField sanitizes a CSV field to prevent CSV injection
func (s *Sanitizer) SanitizeCSVField(field string) string {
	// CSV injection prevention: if field starts with special characters, prefix with tab
	if strings.HasPrefix(field, "=") ||
		strings.HasPrefix(field, "+") ||
		strings.HasPrefix(field, "-") ||
		strings.HasPrefix(field, "@") {
		field = "\t" + field
	}

	// Escape quotes
	field = strings.ReplaceAll(field, `"`, `""`)

	// Wrap in quotes if contains comma, newline, or quote
	if strings.ContainsAny(field, ",\n\"") {
		field = `"` + field + `"`
	}

	return field
}

// SanitizeXML sanitizes XML content
func (s *Sanitizer) SanitizeXML(input string) string {
	// Remove XML declaration if present
	input = regexp.MustCompile(`<\?xml[^>]*\?>`).ReplaceAllString(input, "")

	// Remove DOCTYPE to prevent XXE attacks
	input = regexp.MustCompile(`(?i)<!DOCTYPE[^>]*>`).ReplaceAllString(input, "")

	// Remove comments
	input = regexp.MustCompile(`<!--.*?-->`).ReplaceAllString(input, "")

	return input
}

// SanitizeBase64 sanitizes and validates Base64 input
func (s *Sanitizer) SanitizeBase64(input string) string {
	// Remove whitespace
	input = s.normalizeWhitespace(input)
	input = strings.TrimSpace(input)

	// Remove data URI prefix if present
	if strings.HasPrefix(input, "data:") {
		parts := strings.SplitN(input, ",", 2)
		if len(parts) == 2 {
			input = parts[1]
		}
	}

	return input
}

// SanitizePhoneNumber sanitizes a phone number
func (s *Sanitizer) SanitizePhoneNumber(phone string) string {
	// Remove all non-digit characters except +
	phone = regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Remove leading + if more than one
	phone = regexp.MustCompile(`\++`).ReplaceAllString(phone, "+")

	// Limit length
	if len(phone) > 20 {
		phone = phone[:20]
	}

	return phone
}

// SanitizeCreditCard sanitizes a credit card number (for display only, never store!)
func (s *Sanitizer) SanitizeCreditCard(card string) string {
	// Remove all non-digit characters
	card = regexp.MustCompile(`\D`).ReplaceAllString(card, "")

	// Only keep last 4 digits for display
	if len(card) > 4 {
		card = "****-****-****-" + card[len(card)-4:]
	}

	return card
}

// SanitizeIPAddress sanitizes an IP address
func (s *Sanitizer) SanitizeIPAddress(ip string) string {
	// Trim whitespace
	ip = strings.TrimSpace(ip)

	// Validate IPv4
	if regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(ip) {
		return ip
	}

	// Validate IPv6 (simplified)
	if regexp.MustCompile(`^[\da-fA-F:]+$`).MatchString(ip) {
		return ip
	}

	return ""
}

// SanitizeDomain sanitizes a domain name
func (s *Sanitizer) SanitizeDomain(domain string) string {
	// Trim whitespace and convert to lowercase
	domain = strings.TrimSpace(strings.ToLower(domain))

	// Remove protocol if present
	domain = regexp.MustCompile(`^https?://`).ReplaceAllString(domain, "")

	// Remove path
	domain = strings.Split(domain, "/")[0]

	// Remove port
	domain = strings.Split(domain, ":")[0]

	// Validate domain format
	if !regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`).MatchString(domain) {
		return ""
	}

	return domain
}

// SanitizeSliceAsCSV sanitizes a slice and returns as CSV string
func (s *Sanitizer) SanitizeSliceAsCSV(input []string) string {
	var buf bytes.Buffer
	for i, field := range input {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(s.SanitizeCSVField(field))
	}
	return buf.String()
}

// SanitizeMultiLine sanitizes multi-line text
func (s *Sanitizer) SanitizeMultiLine(input string) string {
	lines := strings.Split(input, "\n")
	sanitized := make([]string, len(lines))

	for i, line := range lines {
		sanitized[i] = s.SanitizeString(line)
	}

	return strings.Join(sanitized, "\n")
}
