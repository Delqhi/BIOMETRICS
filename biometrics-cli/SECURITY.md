# SECURITY.md - Biometrics CLI Security Documentation

**Project:** Biometrics CLI  
**Version:** 1.0.0  
**Date:** February 2026  
**Status:** Sprint 1 Feature Documentation

---

## Table of Contents

1. [Security Overview](#security-overview)
2. [mTLS Configuration Guide](#mtls-configuration-guide)
3. [OAuth2 Provider Setup](#oauth2-provider-setup)
4. [Audit Logging Configuration](#audit-logging-configuration)
5. [Input Validation and Sanitization](#input-validation-and-sanitization)
6. [Security Best Practices 2026](#security-best-practices-2026)
7. [Environment Variables](#environment-variables)
8. [Example Configurations](#example-configurations)

---

## Security Overview

The Biometrics CLI implements enterprise-grade security features as part of Sprint 1. This documentation covers the core security components that protect your application against common vulnerabilities and ensure compliance with security best practices.

### Core Security Features

The security architecture consists of several interconnected components that work together to provide comprehensive protection:

- **Mutual TLS (mTLS):** Provides bidirectional authentication between services, ensuring both client and server verify each other's identity before establishing a connection.
- **OAuth2 Authentication:** Supports multiple identity providers (Google, GitHub, Microsoft) for secure user authentication with standard OAuth2 flows.
- **Audit Logging:** Records all security-relevant events including authentication attempts, authorization decisions, and data access operations.
- **Input Validation:** Implements comprehensive input sanitization to prevent injection attacks, XSS, and other common web vulnerabilities.
- **Token Management:** Secure storage, rotation, and refresh of authentication tokens with encryption at rest.

### Security Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                      BIOMETRICS SECURITY ARCHITECTURE               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐        │
│  │   Client A   │    │   Client B   │    │   Client C   │        │
│  │  (OAuth2)    │    │  (OAuth2)    │    │  (OAuth2)    │        │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘        │
│         │                    │                    │                  │
│         └────────────────────┼────────────────────┘                  │
│                              │                                       │
│                              ▼                                       │
│                   ┌─────────────────────┐                           │
│                   │    mTLS Gateway    │                           │
│                   │  (Biometrics CLI)  │                           │
│                   └──────────┬──────────┘                           │
│                              │                                       │
│         ┌────────────────────┼────────────────────┐                │
│         │                    │                    │                  │
│         ▼                    ▼                    ▼                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │   Audit     │    │  Validation │    │   Cache     │          │
│  │   Logger   │    │   Layer     │    │   (Redis)   │          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## mTLS Configuration Guide

Mutual TLS (mTLS) ensures that both the client and server authenticate each other using certificates. This provides stronger security than traditional TLS where only the server proves its identity.

### Certificate Generation

The Biometrics CLI includes built-in certificate generation capabilities through the `MTLSManager`. Here's how to set it up:

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/delqhi/biometrics/pkg/auth"
)

func main() {
    // Configure mTLS with custom settings
    config := &auth.MTLSConfig{
        CAPath:           "/etc/biometrics/mtls/ca.pem",
        ServerCertPath:   "/etc/biometrics/mtls/server.pem",
        ServerKeyPath:    "/etc/biometrics/mtls/server.key",
        ClientCAPath:     "/etc/biometrics/mtls/client-ca.pem",
        MinVersion:       0x3013, // TLS 1.3
        MaxVersion:       0x3013, // TLS 1.3 only
        ClientAuth:       4,      // RequireAndVerifyClientCert
        RotationInterval: 24 * time.Hour,
        AutoRotate:       true,
    }
    
    manager, err := auth.NewMTLSManager(config)
    if err != nil {
        log.Fatalf("Failed to create MTLS manager: %v", err)
    }
    
    // Get TLS configuration for your server
    tlsConfig, err := manager.GetTLSConfig()
    if err != nil {
        log.Fatalf("Failed to get TLS config: %v", err)
    }
    
    fmt.Printf("TLS Config: %+v\n", tlsConfig)
}
```

### Server Configuration

To configure your HTTP server to use mTLS:

```go
package main

import (
    "crypto/tls"
    "net/http"
    
    "github.com/delqhi/biometrics/pkg/auth"
)

func createMTLSServer(manager *auth.MTLSManager) *http.Server {
    tlsConfig, err := manager.GetTLSConfig()
    if err != nil {
        panic(err)
    }
    
    server := &http.Server{
        Addr:      ":8443",
        TLSConfig: tlsConfig,
        Handler:   yourHandler(),
    }
    
    return server
}
```

### Client Certificate Generation

For service-to-service communication, generate client certificates:

```go
// Generate a client certificate for a specific service
clientCert, err := manager.GenerateClientCertificate(
    "user-service",           // Client identifier
    90 * 24 * time.Hour,     // Validity duration (90 days)
)
if err != nil {
    return err
}

// Use the certificate in your client
client := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            Certificates: []tls.Certificate{*clientCert},
            RootCAs:      certPool,
        },
    },
}
```

### Certificate Validation

The mTLS manager provides built-in certificate validation:

```go
// Validate a client certificate
func validateClientCert(manager *auth.MTLSManager, certPath string) error {
    info, err := manager.GetCertificateInfo(certPath)
    if err != nil {
        return fmt.Errorf("failed to get certificate info: %w", err)
    }
    
    // Check expiration
    daysLeft := info.NotAfter.Sub(time.Now()).Hours() / 24
    if daysLeft < 7 {
        fmt.Printf("WARNING: Certificate expires in %.0f days\n", daysLeft)
    }
    
    // Check if rotation is needed
    if manager.NeedsRotation() {
        fmt.Println("WARNING: Certificate rotation needed")
    }
    
    return nil
}
```

### Certificate Rotation

Configure automatic certificate rotation:

```go
config := &auth.MTLSConfig{
    // ... other config
    RotationInterval: 24 * 7 * time.Hour, // Weekly rotation
    AutoRotate:       true,
}

manager, err := auth.NewMTLSManager(config)
if err != nil {
    log.Fatal(err)
}

// Manual rotation when needed
if err := manager.RotateCertificates(); err != nil {
    log.Printf("Rotation failed: %v", err)
}
```

---

## OAuth2 Provider Setup

Biometrics CLI supports multiple OAuth2 providers for user authentication. This section covers setup for Google, GitHub, and Microsoft Azure AD.

### Google OAuth2 Setup

#### Step 1: Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create a new project or select existing
3. Navigate to APIs & Services > Credentials
4. Create OAuth 2.0 Client ID credentials
5. Set authorized redirect URIs to `http://localhost:8080/auth/callback/google`

#### Step 2: Configure in Code

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/delqhi/biometrics/pkg/auth/oauth2"
)

func main() {
    client, err := oauth2.NewOAuth2Client([]oauth2.OAuth2Config{
        {
            ClientID:     "your-google-client-id",
            ClientSecret: "your-google-client-secret",
            RedirectURL:  "http://localhost:8080/auth/callback/google",
            Provider:     oauth2.ProviderGoogle,
            Scopes: []string{
                "openid",
                "https://www.googleapis.com/auth/userinfo.email",
                "https://www.googleapis.com/auth/userinfo.profile",
            },
        },
    })
    if err != nil {
        panic(err)
    }
    
    // Generate authorization URL
    authURL, state, err := client.GetAuthURL(oauth2.ProviderGoogle)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Authorization URL:", authURL)
    fmt.Println("State (save this):", state)
}
```

### GitHub OAuth2 Setup

#### Step 1: Register OAuth App

1. Go to GitHub Settings > Developer Settings > OAuth Apps
2. Click "New OAuth App"
3. Set Homepage URL to `http://localhost:8080`
4. Set Authorization callback URL to `http://localhost:8080/auth/callback/github`

#### Step 2: Configure in Code

```go
client, err := oauth2.NewOAuth2Client([]oauth2.OAuth2Config{
    {
        ClientID:     "your-github-client-id",
        ClientSecret: "your-github-client-secret",
        RedirectURL:  "http://localhost:8080/auth/callback/github",
        Provider:     oauth2.ProviderGitHub,
        Scopes: []string{
            "read:user",
            "user:email",
            "repo",
        },
    },
})
```

### Microsoft Azure AD Setup

```go
// For Azure AD, use the generic OAuth2 config with Microsoft endpoints
client, err := oauth2.NewOAuth2Client([]oauth2.OAuth2Config{
    {
        ClientID:     "your-azure-client-id",
        ClientSecret: "your-azure-client-secret",
        RedirectURL:  "http://localhost:8080/auth/callback/azure",
        Provider:     "azure",  // Custom provider
        Scopes: []string{
            "openid",
            "profile",
            "email",
            "User.Read",
        },
    },
})
```

### OAuth2 Authentication Flow

```go
// Complete authentication flow
func handleOAuthCallback(client *oauth2.OAuth2Client, provider oauth2.ProviderType, callbackURL, expectedState string) (*oauth2.UserInfo, error) {
    // Complete the authentication
    userInfo, err := client.CompleteAuth(context.Background(), provider, callbackURL, expectedState)
    if err != nil {
        return nil, fmt.Errorf("auth failed: %w", err)
    }
    
    // UserInfo contains:
    // - ID: Unique user ID from provider
    // - Email: User's email address
    // - Name: Full name
    // - Picture: Profile picture URL
    // - Provider: Provider name
    
    fmt.Printf("Authenticated: %s <%s>\n", userInfo.Name, userInfo.Email)
    
    return userInfo, nil
}
```

### Token Management

```go
// Get stored token for a user
tokenManager := client.GetTokenManager()
token := tokenManager.Get("google", "user-id-123")

// Refresh token if expired
if token != nil && !token.Valid() {
    newToken, err := client.RefreshToken(context.Background(), "google", token.RefreshToken)
    if err != nil {
        return err
    }
    // Store new token
    tokenManager.Store("google", "user-id-123", newToken)
}

// Logout and revoke token
err = client.Logout(context.Background(), oauth2.ProviderGoogle, "user-id-123")
```

---

## Audit Logging Configuration

Security audit logging is critical for compliance and incident response. Biometrics CLI provides comprehensive audit logging through the `audit` package.

### Basic Configuration

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/delqhi/biometrics/pkg/audit"
)

func main() {
    config := &audit.AuditConfig{
        StoragePath:        "/var/log/biometrics/audit",
        StorageType:        audit.StorageTypeFile,
        MaxSize:            100 * 1024 * 1024, // 100MB
        RetentionDays:      90,
        EnableCompression:  true,
        EnableEncryption:   false,
        FlushInterval:      5 * time.Second,
        QueueSize:          1000,
    }
    
    auditor, err := audit.NewAuditor(config)
    if err != nil {
        panic(err)
    }
    defer auditor.Stop()
    
    // Log authentication success
    err = auditor.LogAuthentication("user-123", true, "oauth2", "192.168.1.100")
    if err != nil {
        fmt.Printf("Failed to log: %v\n", err)
    }
}
```

### Logging Security Events

```go
// Log various security events
func logSecurityEvents(auditor *audit.Auditor) {
    // Authentication success
    auditor.LogAuthentication("user-123", true, "oauth2", "192.168.1.100")
    
    // Authentication failure
    auditor.LogAuthentication("unknown-user", false, "password", "10.0.0.55")
    
    // Authorization granted
    auditor.LogAuthorization("user-123", "/api/admin", "DELETE", true)
    
    // Authorization denied
    auditor.LogAuthorization("user-456", "/api/admin", "DELETE", false)
    
    // Data access
    auditor.LogDataAccess("user-123", "medical-records", "READ", "rec-789")
    
    // Security alerts
    auditor.LogSecurityEvent(
        audit.EventSecurityAlert,
        "system",
        "Multiple failed login attempts detected",
        "high",
    )
    
    // Configuration changes
    auditor.Log(
        audit.EventConfigChange,
        "admin-user",
        "config-change",
        "security-settings",
        map[string]interface{}{
            "key":        "max_login_attempts",
            "old_value": "5",
            "new_value": "3",
        },
    )
}
```

### Querying Audit Logs

```go
// Query audit logs
func queryAuditLogs(auditor *audit.Auditor) error {
    query := &audit.AuditQuery{
        StartTime:  time.Now().Add(-24 * time.Hour),
        EndTime:    time.Now(),
        EventTypes: []audit.EventType{audit.EventAuthFailure},
        Limit:      100,
        SortBy:     "timestamp",
        SortOrder:  "desc",
    }
    
    results, err := auditor.Query(context.Background(), query)
    if err != nil {
        return err
    }
    
    fmt.Printf("Found %d events\n", results.TotalCount)
    for _, event := range results.Events {
        fmt.Printf("[%s] %s: %s\n", 
            event.Timestamp.Format(time.RFC3339),
            event.EventType,
            event.Actor,
        )
    }
    
    return nil
}
```

### Export Audit Logs

```go
// Export to JSON
jsonData, err := auditor.Export(
    time.Now().Add(-30 * 24 * time.Hour),
    time.Now(),
    audit.ExportFormatJSON,
)

// Export to CSV
csvData, err := auditor.Export(
    time.Now().Add(-30 * 24 * time.Hour),
    time.Now(),
    audit.ExportFormatCSV,
)
```

### Audit Statistics

```go
// Get audit statistics
stats, err := auditor.GetStats()
if err != nil {
    return err
}

fmt.Printf("Total Events: %d\n", stats.TotalEvents)
fmt.Printf("Events by Type: %+v\n", stats.EventsByType)
fmt.Printf("Events by Actor: %+v\n", stats.EventsByActor)
fmt.Printf("Storage Size: %d bytes\n", stats.StorageSize)
fmt.Printf("Avg Events/Day: %.2f\n", stats.AvgEventsPerDay)
```

---

## Input Validation and Sanitization

The validation package provides comprehensive input validation and sanitization to protect against injection attacks, XSS, and other vulnerabilities.

### Basic Validation

```go
package main

import (
    "fmt"
    
    "github.com/delqhi/biometrics/pkg/validation"
)

func main() {
    // Validate email
    email := "user@example.com"
    if err := validation.ValidateEmail(email); err != nil {
        fmt.Printf("Invalid email: %v\n", err)
    }
    
    // Validate URL
    url := "https://example.com/path"
    if err := validation.ValidateURL(url); err != nil {
        fmt.Printf("Invalid URL: %v\n", err)
    }
    
    // Validate required fields
    data := map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
    }
    
    rules := validation.Rules{
        "name":  validation.Required,
        "email": validation.Required | validation.IsEmail,
    }
    
    if err := validation.Validate(data, rules); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    }
}
```

### Input Sanitization

```go
// Sanitize user input to prevent XSS
func sanitizeUserInput(input string) string {
    // HTML escape
    sanitized := validation.SanitizeHTML(input)
    
    // Remove dangerous protocols
    sanitized = validation.StripDangerousURLs(sanitized)
    
    // Normalize Unicode
    sanitized = validation.NormalizeUnicode(sanitized)
    
    return sanitized
}

// Sanitize for SQL (parameterized queries recommended)
func sanitizeForSQL(input string) string {
    // This is for logging/display - always use parameterized queries!
    return validation.EscapeSpecialChars(input)
}

// Validate and sanitize file paths
func validateFilePath(path string) error {
    // Check for path traversal
    if validation.ContainsPathTraversal(path) {
        return fmt.Errorf("path traversal detected")
    }
    
    // Validate against allowed patterns
    if !validation.MatchesPattern(path, `^/data/[a-z0-9/]+$`) {
        return fmt.Errorf("invalid path format")
    }
    
    return nil
}
```

### Custom Validation Rules

```go
// Create custom validation rules
func customValidation() {
    rules := validation.Rules{
        "username": validation.Required | validation.MinLength(3) | validation.MaxLength(20),
        "password": validation.Required | validation.MinLength(8) | validation.StrongPassword,
        "age":      validation.Required | validation.Min(0) | validation.Max(150),
        "phone":    validation.Pattern(`^\+?[1-9]\d{1,14}$`),
    }
    
    data := map[string]interface{}{
        "username": "john_doe",
        "password": "SecureP@ss123",
        "age":      25,
        "phone":    "+1234567890",
    }
    
    if err := validation.Validate(data, rules); err != nil {
        fmt.Printf("Validation error: %v\n", err)
    }
}
```

### Middleware Integration

```go
import "github.com/delqhi/biometrics/pkg/validation"

// HTTP middleware for request validation
func validationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate request body
        if r.Body != nil {
            // Parse and validate
            var data map[string]interface{}
            if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
            }
            
            // Validate against rules
            if err := validation.Validate(data, getValidationRules(r.URL.Path)); err != nil {
                http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
                return
            }
        }
        
        next.ServeHTTP(w, r)
    })
}

func getValidationRules(path string) validation.Rules {
    switch {
    case path == "/api/users":
        return validation.Rules{
            "email":    validation.Required | validation.IsEmail,
            "name":     validation.Required | validation.MinLength(1),
            "password": validation.Required | validation.MinLength(8),
        }
    case path == "/api/login":
        return validation.Rules{
            "email":    validation.Required | validation.IsEmail,
            "password": validation.Required,
        }
    default:
        return validation.Rules{}
    }
}
```

---

## Security Best Practices 2026

The following best practices ensure your Biometrics CLI deployment meets enterprise security standards.

### 1. Certificate Management

- **Use TLS 1.3 only:** Disable older protocol versions
- **Implement certificate rotation:** Auto-rotate every 24-90 days
- **Use strong key sizes:** Minimum 2048-bit RSA, 256-bit ECC
- **Store keys securely:** Use HashiCorp Vault or similar

```go
config := &auth.MTLSConfig{
    MinVersion: tls.VersionTLS13,
    MaxVersion: tls.VersionTLS13,
    CipherSuites: []uint16{
        tls.TLS_AES_256_GCM_SHA384,
        tls.TLS_CHACHA20_POLY1305_SHA256,
        tls.TLS_AES_128_GCM_SHA256,
    },
}
```

### 2. OAuth2 Security

- **Use PKCE:** Implement Proof Key for Code Exchange
- **Validate state parameter:** Prevent CSRF attacks
- **Use short-lived tokens:** Access tokens should expire within 1 hour
- **Store tokens encrypted:** Never store tokens in plain text

```go
// Always use state parameter
authURL, state, err := client.GetAuthURL(provider)
// Store state in session/state store
session.Set("oauth_state", state)

// Validate on callback
if err := client.ValidateState(callbackState, storedState); err != nil {
    return err // State mismatch - potential CSRF
}
```

### 3. Audit Logging

- **Log all authentication events:** Success and failure
- **Log authorization decisions:** Access grants and denials
- **Log data access:** Especially sensitive data
- **Implement log rotation:** Prevent disk exhaustion

### 4. Input Validation

- **Validate on server:** Never trust client-side validation
- **Use allowlists:** Prefer allowed characters over blocked
- **Sanitize outputs:** Escape data before display
- **Parameterize queries:** Never concatenate user input into SQL

### 5. Secrets Management

- **Use environment variables:** Never hardcode secrets
- **Rotate secrets regularly:** Implement automated rotation
- **Use secret scanning:** Detect accidental commits

---

## Environment Variables

The following environment variables configure security features:

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `BIOMETRICS_MTLS_CA_PATH` | Path to CA certificate | No | `/tmp/biometrics/mtls/ca.pem` |
| `BIOMETRICS_MTLS_CERT_PATH` | Path to server certificate | No | `/tmp/biometrics/mtls/server.pem` |
| `BIOMETRICS_MTLS_KEY_PATH` | Path to server private key | No | `/tmp/biometrics/mtls/server.key` |
| `BIOMETRICS_MTLS_AUTO_ROTATE` | Enable auto certificate rotation | No | `false` |
| `BIOMETRICS_OAUTH_GOOGLE_CLIENT_ID` | Google OAuth client ID | No | - |
| `BIOMETRICS_OAUTH_GOOGLE_CLIENT_SECRET` | Google OAuth client secret | No | - |
| `BIOMETRICS_OAUTH_GITHUB_CLIENT_ID` | GitHub OAuth client ID | No | - |
| `BIOMETRICS_OAUTH_GITHUB_CLIENT_SECRET` | GitHub OAuth client secret | No | - |
| `BIOMETRICS_AUDIT_PATH` | Path for audit logs | No | `/tmp/biometrics/audit` |
| `BIOMETRICS_AUDIT_RETENTION_DAYS` | Days to retain audit logs | No | `90` |
| `BIOMETRICS_JWT_SECRET` | Secret for JWT signing | Yes | - |
| `BIOMETRICS_ENCRYPTION_KEY` | Key for token encryption | Yes | - |

---

## Example Configurations

### Production Configuration

```yaml
# config.yaml
security:
  mtls:
    enabled: true
    ca_path: /etc/biometrics/mtls/ca.pem
    server_cert_path: /etc/biometrics/mtls/server.pem
    server_key_path: /etc/biometrics/mtls/server.key
    auto_rotate: true
    rotation_interval: 24h
    
  oauth:
    providers:
      google:
        client_id: ${BIOMETRICS_OAUTH_GOOGLE_CLIENT_ID}
        client_secret: ${BIOMETRICS_OAUTH_GOOGLE_CLIENT_SECRET}
        redirect_url: https://biometrics.example.com/auth/callback/google
      github:
        client_id: ${BIOMETRICS_OAUTH_GITHUB_CLIENT_ID}
        client_secret: ${BIOMETRICS_OAUTH_GITHUB_CLIENT_SECRET}
        redirect_url: https://biometrics.example.com/auth/callback/github
        
  audit:
    enabled: true
    storage_path: /var/log/biometrics/audit
    storage_type: file
    retention_days: 365
    compression: true
    
  validation:
    strict_mode: true
    max_request_size: 1048576
    allowed_content_types:
      - application/json
      - multipart/form-data
```

### Development Configuration

```yaml
# config-dev.yaml
security:
  mtls:
    enabled: false
    
  oauth:
    providers:
      google:
        client_id: dev-google-client-id
        client_secret: dev-google-secret
        redirect_url: http://localhost:8080/auth/callback/google
        
  audit:
    enabled: true
    storage_path: /tmp/biometrics/audit-dev
    storage_type: memory
    retention_days: 7
```

---

## Related Documentation

- [Performance Documentation](./PERFORMANCE.md)
- [API Reference](./docs/api/)
- [Deployment Guide](./docs/deployment/)
- [Troubleshooting](./docs/troubleshooting/)

---

*Document Version: 1.0.0*  
*Last Updated: February 2026*  
*Compliant with Enterprise Practices Feb 2026*
