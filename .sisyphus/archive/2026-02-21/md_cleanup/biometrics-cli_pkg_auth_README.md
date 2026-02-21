# Authentication Package

**Purpose:** Authentication and authorization utilities for biometrics CLI

## Overview

This package provides authentication and authorization functionality for the biometrics CLI, supporting multiple authentication providers and credential management.

## Supported Providers

| Provider | Type | Features |
|----------|------|----------|
| Azure AD | OAuth2/OIDC | SSO, MFA |
| Google | OAuth2 | SSO |
| GitHub | OAuth | Teams |
| Local | Credentials | API keys |

## Installation

```go
import "github.com/delqhi/biometrics/pkg/auth"
```

## Usage

### Initializing Provider

```go
// Create Azure AD provider
provider, err := auth.NewProvider(ctx, auth.ProviderConfig{
    Type: "azure",
    TenantID: "your-tenant-id",
    ClientID: "your-client-id",
})
```

### Authentication Flow

```go
// Start authentication
url, err := provider.AuthURL()
// Redirect user to URL

// Exchange code for token
token, err := provider.Exchange(ctx, code)

// Use token
client := provider.Client(ctx, token)
```

### Token Management

```go
// Store token
err := tokenManager.Store(ctx, "provider", token)

// Retrieve token
token, err := tokenManager.Get(ctx, "provider")

// Refresh token
newToken, err := token.Refresh(ctx)

// Revoke token
err := token.Revoke(ctx)
```

## Configuration

### Azure AD

```yaml
auth:
  providers:
    azure:
      type: oauth2
      tenant_id: "${AZURE_TENANT_ID}"
      client_id: "${AZURE_CLIENT_ID}"
      client_secret: "${AZURE_CLIENT_SECRET}"
      scopes:
        - openid
        - profile
        - email
```

### GitHub

```yaml
auth:
  providers:
    github:
      type: oauth2
      client_id: "${GITHUB_CLIENT_ID}"
      client_secret: "${GITHUB_CLIENT_SECRET}"
      scopes:
        - read:user
        - repo
```

## Security Features

### Token Encryption

```go
// Encrypt token at rest
encrypted, err := auth.EncryptToken(token, encryptionKey)

// Decrypt token
token, err := auth.DecryptToken(encrypted, encryptionKey)
```

### MFA Support

```go
// Handle MFA challenge
challenge, err := provider.MFAChallenge(ctx, token)
// Support TOTP, SMS, Email
```

## Error Handling

```go
_, err := provider.Authenticate(ctx, creds)
if err != nil {
    switch {
    case errors.Is(err, auth.ErrInvalidCredentials):
        // Show login form again
    case errors.Is(err, auth.ErrMFARequired):
        // Prompt for MFA code
    case errors.Is(err, auth.ErrTokenExpired):
        // Refresh or re-authenticate
    }
}
```

## Testing

```bash
# Run tests
go test ./pkg/auth/... -v

# Run with coverage
go test ./pkg/auth/... -coverprofile=coverage.out

# Run integration tests
go test ./pkg/auth/... -tags=integration -v
```

## Related Documentation

- [Security Configuration](../docs/security/)
- [Provider Setup](../config/providers.md)
- [OAuth Flow](../docs/oauth/)
