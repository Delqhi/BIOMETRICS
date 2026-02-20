# OAuth2 Package - Comprehensive Documentation

**Package:** github.com/delqhi/biometrics/pkg/auth/oauth2  
**Version:** 1.0.0  
**Date:** February 2026  
**Status:** Sprint 1 Feature

---

## Table of Contents

1. [Package Overview](#package-overview)
2. [Supported Providers](#supported-providers)
3. [Installation](#installation)
4. [Quick Start per Provider](#quick-start-per-provider)
5. [API Reference](#api-reference)
6. [Token Management](#token-management)
7. [Security Considerations](#security-considerations)

---

## Package Overview

The oauth2 package provides a complete OAuth2 client implementation with support for multiple identity providers. It handles the complete OAuth2 flow including authorization, token exchange, refresh, and revocation.

### Key Features

- **Multi-Provider Support:** Google, GitHub, and custom OAuth2 providers
- **Complete OAuth2 Flow:** Authorization code flow with PKCE support
- **Token Management:** Automatic token refresh and secure storage
- **User Info Retrieval:** Unified user profile interface across providers
- **Token Validation:** Verify tokens without API calls
- **Token Revocation:** Properly revoke tokens on logout
- **State Validation:** CSRF protection with state parameter
- **Thread-Safe:** Concurrent request handling

### Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                       OAUTH2 PACKAGE ARCHITECTURE                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                    OAuth2Client                              │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │   │
│  │  │   Config   │  │    Token   │  │   HTTP     │         │   │
│  │  │  Manager   │  │   Manager   │  │   Client   │         │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘         │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                              │                                       │
│         ┌────────────────────┼────────────────────┐                │
│         │                    │                    │                  │
│         ▼                    ▼                    ▼                  │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐          │
│  │   Google   │      │   GitHub   │      │  Custom    │          │
│  │  Provider  │      │  Provider  │      │  Provider  │          │
│  └─────────────┘      └─────────────┘      └─────────────┘          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### OAuth2 Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                        OAUTH2 FLOW DIAGRAM                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   User         Client App         Auth Server        API Server     │
│    │               │                  │                   │          │
│    │  1.Click     │                  │                   │          │
│    │──────────────>│                  │                   │          │
│    │               │                  │                   │          │
│    │  2.Redirect   │                  │                   │          │
│    │  (Auth URL)   │                  │                   │          │
│    │<──────────────│                  │                   │          │
│    │               │                  │                   │          │
│    │  3.Login &   │                  │                   │          │
│    │  Consent     │                  │                   │          │
│    │──────────────>│───────────────>│                   │          │
│    │               │                  │                   │          │
│    │  4.Callback  │                  │                   │          │
│    │  (Code)      │                  │                   │          │
│    │<──────────────│<───────────────│                   │          │
│    │               │                  │                   │          │
│    │               │  5.Exchange     │                   │          │
│    │               │  (Code for      │                   │          │
│    │               │   Token)        │                   │          │
│    │               │───────────────>│                   │          │
│    │               │<───────────────│                   │          │
│    │               │                  │                   │          │
│    │               │  6.API Request  │                   │          │
│    │               │  (with Token)  │                   │          │
│    │               │───────────────────────────────────>│          │
│    │               │<───────────────────────────────────│          │
│    │               │                  │                   │          │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Supported Providers

### Google

- **OAuth Version:** OAuth 2.0
- **Token Endpoint:** https://oauth2.googleapis.com/token
- **Authorization Endpoint:** https://accounts.google.com/o/oauth2/v2/auth
- **User Info Endpoint:** https://www.googleapis.com/oauth2/v2/userinfo
- **Token Revocation:** https://oauth2.googleapis.com/revoke
- **Required Scopes:** `openid`, `email`, `profile`

### GitHub

- **OAuth Version:** OAuth 2.0
- **Token Endpoint:** https://github.com/login/oauth/access_token
- **Authorization Endpoint:** https://github.com/login/oauth/authorize
- **User Info Endpoint:** https://api.github.com/user
- **Token Revocation:** Not supported by GitHub
- **Required Scopes:** `read:user`, `user:email`

### Custom Providers

The package supports custom OAuth2 providers by configuration:

```go
config := oauth2.OAuth2Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    RedirectURL:  "https://your-app.com/callback",
    Provider:     "custom-provider",
    Scopes:       []string{"openid", "profile"},
}
```

---

## Installation

### Prerequisites

- Go 1.21 or later
- OAuth2 credentials from your provider

### Installation Command

```bash
go get github.com/delqhi/biometrics/pkg/auth/oauth2
```

### Dependencies

```go
import (
    "context"           // Standard library
    "fmt"               // Standard library
    "net/http"          // Standard library
    "sync"              // Standard library
    "time"              // Standard library
    
    "golang.org/x/oauth2"  // OAuth2 support
)
```

---

## Quick Start per Provider

### Google Quick Start

#### Step 1: Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create new project or select existing
3. Navigate to APIs & Services > Credentials
4. Create OAuth 2.0 Client ID
5. Set redirect URIs

#### Step 2: Configure the Client

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/delqhi/biometrics/pkg/auth/oauth2"
)

func main() {
    // Create OAuth2 client with Google configuration
    client, err := oauth2.NewOAuth2Client([]oauth2.OAuth2Config{
        {
            ClientID:     "YOUR-GOOGLE-CLIENT-ID",
            ClientSecret: "YOUR-GOOGLE-CLIENT-SECRET",
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
        log.Fatal(err)
    }
    
    // Get authorization URL
    authURL, state, err := client.GetAuthURL(oauth2.ProviderGoogle)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Authorization URL:", authURL)
    fmt.Println("State (save this):", state)
    
    // After user authorizes, exchange code for token
    // In practice: retrieve code from callback URL
    ctx := context.Background()
    token, err := client.Exchange(ctx, oauth2.ProviderGoogle, "AUTHORIZATION_CODE")
    if err != nil {
        log.Fatal(err)
    }
    
    // Get user info
    userInfo, err := client.GetUserInfo(ctx, oauth2.ProviderGoogle, token)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User: %s <%s>\n", userInfo.Name, userInfo.Email)
}
```

### GitHub Quick Start

#### Step 1: Register OAuth App

1. Go to GitHub Settings > Developer Settings > OAuth Apps
2. Click "New OAuth App"
3. Fill in application details
4. Note the Client ID and Client Secret

#### Step 2: Configure the Client

```go
func githubExample() {
    client, err := oauth2.NewOAuth2Client([]oauth2.OAuth2Config{
        {
            ClientID:     "YOUR-GITHUB-CLIENT-ID",
            ClientSecret: "YOUR-GITHUB-CLIENT-SECRET",
            RedirectURL:  "http://localhost:8080/auth/callback/github",
            Provider:     oauth2.ProviderGitHub,
            Scopes: []string{
                "read:user",
                "user:email",
            },
        },
    })
    if err != nil {
        panic(err)
    }
    
    // Get authorization URL
    authURL, state, err := client.GetAuthURL(oauth2.ProviderGitHub)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("GitHub Authorization:", authURL)
    fmt.Println("State:", state)
}
```

### Complete Web Flow

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/delqhi/biometrics/pkg/auth/oauth2"
)

type AuthHandler struct {
    client  *oauth2.OAuth2Client
    states  map[string]time.Time  // State -> Created time
}

func NewAuthHandler(client *oauth2.OAuth2Client) *AuthHandler {
    return &AuthHandler{
        client: client,
        states: make(map[string]time.Time),
    }
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
    provider := r.URL.Query().Get("provider")
    if provider == "" {
        provider = string(oauth2.ProviderGoogle)
    }
    
    // Get authorization URL
    authURL, state, err := h.client.GetAuthURL(oauth2.ProviderType(provider))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Store state with expiration (10 minutes)
    h.states[state] = time.Now().Add(10 * time.Minute)
    
    // Redirect user to provider
    http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
    // Get code and state from callback
    code := r.URL.Query().Get("code")
    state := r.URL.Query().Get("state")
    
    // Validate state
    if createdAt, exists := h.states[state]; !exists || time.Now().After(createdAt) {
        http.Error(w, "Invalid or expired state", http.StatusBadRequest)
        return
    }
    delete(h.states, state)
    
    // Complete authentication
    userInfo, err := h.client.CompleteAuth(
        context.Background(),
        oauth2.ProviderGoogle,
        r.URL.String(),
        state,
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "Welcome %s!", userInfo.Name)
}
```

---

## API Reference

### Types

#### OAuth2Client

Main OAuth2 client that manages authentication:

```go
type OAuth2Client struct {
    // Contains private fields
}
```

#### OAuth2Config

Configuration for an OAuth2 provider:

```go
type OAuth2Config struct {
    ClientID     string        // OAuth2 client ID
    ClientSecret string        // OAuth2 client secret
    RedirectURL  string        // OAuth2 redirect URI
    Scopes       []string      // Requested scopes
    Provider     ProviderType  // Provider identifier
}
```

#### UserInfo

User information from OAuth2 provider:

```go
type UserInfo struct {
    ID          string `json:"id"`           // Provider-specific user ID
    Email       string `json:"email"`        // User email
    Name        string `json:"name"`         // Full name
    Picture     string `json:"picture"`      // Profile picture URL
    Provider    string `json:"provider"`     // Provider name
    AccessToken string `json:"access_token"` // OAuth2 access token
}
```

#### ProviderType

Supported OAuth2 providers:

```go
const (
    ProviderGoogle ProviderType = "google"
    ProviderGitHub ProviderType = "github"
)
```

### Functions

#### NewOAuth2Client

Create a new OAuth2 client with one or more provider configurations:

```go
func NewOAuth2Client(configs []OAuth2Config) (*OAuth2Client, error)
```

### Methods

#### GetAuthURL

Generate authorization URL for user redirection:

```go
func (c *OAuth2Client) GetAuthURL(provider ProviderType) (authURL, state string, error)
```

#### Exchange

Exchange authorization code for access token:

```go
func (c *OAuth2Client) Exchange(ctx context.Context, provider ProviderType, code string) (*oauth2.Token, error)
```

#### GetUserInfo

Retrieve user information from provider:

```go
func (c *OAuth2Client) GetUserInfo(ctx context.Context, provider ProviderType, token *oauth2.Token) (*UserInfo, error)
```

#### RefreshToken

Refresh an expired access token:

```go
func (c *OAuth2Client) RefreshToken(ctx context.Context, provider ProviderType, refreshToken string) (*oauth2.Token, error)
```

#### RevokeToken

Revoke an access or refresh token:

```go
func (c *OAuth2Client) RevokeToken(ctx context.Context, provider ProviderType, token string) error
```

#### ValidateToken

Validate an access token:

```go
func (c *OAuth2Client) ValidateToken(ctx context.Context, provider ProviderType, token string) (bool, error)
```

#### ValidateState

Validate OAuth2 state parameter for CSRF protection:

```go
func (c *OAuth2Client) ValidateState(state, expected string) error
```

#### CompleteAuth

Complete the full OAuth2 flow in one call:

```go
func (c *OAuth2Client) CompleteAuth(ctx context.Context, provider ProviderType, callbackURL, expectedState string) (*UserInfo, error)
```

#### Logout

Logout user and revoke token:

```go
func (c *OAuth2Client) Logout(ctx context.Context, provider ProviderType, userID string) error
```

#### ListProviders

List configured OAuth2 providers:

```go
func (c *OAuth2Client) ListProviders() []ProviderType
```

#### IsConfigured

Check if a provider is configured:

```go
func (c *OAuth2Client) IsConfigured(provider ProviderType) bool
```

---

## Token Management

### Token Storage

The TokenManager handles secure token storage:

```go
// Get token manager
tokenManager := client.GetTokenManager()

// Store token
tokenManager.Store("google", "user-123", token)

// Retrieve token
token := tokenManager.Get("google", "user-123")

// Delete token
tokenManager.Delete("google", "user-123")
```

### Automatic Token Refresh

Tokens are automatically refreshed when expired:

```go
// TokenManager handles refresh automatically
token, err := client.Exchange(ctx, provider, code)
if err != nil {
    return err
}

// Store token (TokenManager will handle refresh)
tokenManager.Store(string(provider), userID, token)

// Later: get token (automatically refreshed if expired)
token = tokenManager.Get(string(provider), userID)
```

### Manual Token Refresh

```go
// Manual refresh when needed
func refreshToken(client *oauth2.OAuth2Client, provider oauth2.ProviderType, userID string) error {
    tokenManager := client.GetTokenManager()
    currentToken := tokenManager.Get(string(provider), userID)
    
    if currentToken == nil || !currentToken.Valid() {
        // Need to refresh
        newToken, err := client.RefreshToken(
            context.Background(),
            provider,
            currentToken.RefreshToken,
        )
        if err != nil {
            return err
        }
        
        // Store new token
        tokenManager.Store(string(provider), userID, newToken)
    }
    
    return nil
}
```

### Token Validation

```go
// Validate token without API call
func validateToken(client *oauth2.OAuth2Client, provider oauth2.ProviderType, tokenString string) (bool, error) {
    return client.ValidateToken(context.Background(), provider, tokenString)
}
```

---

## Security Considerations

### State Parameter Validation

Always validate the state parameter to prevent CSRF attacks:

```go
// Generate state and store it
authURL, state, err := client.GetAuthURL(provider)

// IMPORTANT: Store state in session or temporary storage
session.Set("oauth_state", state)
session.Set("oauth_state_created", time.Now())

// On callback: validate state
callbackState := r.URL.Query().Get("state")
storedState := session.Get("oauth_state")

if err := client.ValidateState(callbackState, storedState); err != nil {
    // State mismatch - possible CSRF attack
    return err
}
```

### Secure Token Storage

Never store tokens in plain text:

```go
// Good: Use encrypted storage
tokenManager := NewTokenManager(TokenManagerOptions{
    EncryptionKey: encryptionKey,  // AES-256 key
})

// Token is encrypted at rest
tokenManager.Store(provider, userID, token)
```

### Token Expiration

Check token expiration:

```go
func isTokenExpired(token *oauth2.Token) bool {
    return token.Expiry.Before(time.Now())
}

// Use token if valid
if !isTokenExpired(token) {
    // Use token
} else if token.RefreshToken != "" {
    // Refresh token
} else {
    // Re-authenticate
}
```

### HTTPS Requirement

Always use HTTPS in production:

```go
// Configure redirect URL with HTTPS
config := oauth2.OAuth2Config{
    // ...
    RedirectURL: "https://your-domain.com/auth/callback",
}
```

### Scope Minimization

Request only necessary scopes:

```go
// Good: Minimal scopes
Scopes: []string{
    "openid",
    "email",
}

// Avoid: Excessive scopes
Scopes: []string{
    "openid",
    "email",
    "profile",
    "https://www.googleapis.com/auth/calendar",  // Not needed
    "https://www.googleapis.com/auth/drive",      // Not needed
}
```

### Error Handling

Handle OAuth2 errors properly:

```go
func handleOAuthError(err error) string {
    switch {
    case strings.Contains(err.Error(), "invalid_grant"):
        return "Authorization code expired or already used"
    case strings.Contains(err.Error(), "invalid_client"):
        return "Invalid client credentials"
    case strings.Contains(err.Error(), "unauthorized"):
        return "Unauthorized - check permissions"
    default:
        return fmt.Sprintf("OAuth2 error: %v", err)
    }
}
```

---

## Example: Complete Integration

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/delqhi/biometrics/pkg/auth/oauth2"
)

func main() {
    // 1. Create OAuth2 client
    client, err := oauth2.NewOAuth2Client([]oauth2.OAuth2Config{
        {
            ClientID:     "client-id",
            ClientSecret: "client-secret",
            RedirectURL:  "http://localhost:8080/callback",
            Provider:     oauth2.ProviderGoogle,
            Scopes: []string{
                "openid",
                "email",
                "profile",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. Get authorization URL (for login page)
    authURL, state, err := client.GetAuthURL(oauth2.ProviderGoogle)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Login at:", authURL)
    fmt.Println("State:", state)
    
    // 3. Simulate callback (in real app: from HTTP request)
    callbackURL := fmt.Sprintf("http://localhost:8080/callback?code=test-code&state=%s", state)
    
    // 4. Complete authentication
    userInfo, err := client.CompleteAuth(
        context.Background(),
        oauth2.ProviderGoogle,
        callbackURL,
        state,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // 5. Use user info
    fmt.Printf("Logged in as: %s (%s)\n", userInfo.Name, userInfo.Email)
    fmt.Printf("User ID: %s\n", userInfo.ID)
    fmt.Printf("Provider: %s\n", userInfo.Provider)
    
    // 6. Logout (optional)
    if err := client.Logout(context.Background(), oauth2.ProviderGoogle, userInfo.ID); err != nil {
        log.Printf("Logout warning: %v", err)
    }
}
```

---

## Related Documentation

- [Security Configuration](../SECURITY.md)
- [Auth Package Overview](../pkg/auth/README.md)
- [Best Practices](https://tools.ietf.org/html/rfc8252)

---

*Document Version: 1.0.0*  
*Last Updated: February 2026*  
*Compliant with Enterprise Practices Feb 2026*
