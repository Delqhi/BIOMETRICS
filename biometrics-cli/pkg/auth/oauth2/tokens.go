package oauth2

import (
	"sync"
	"time"

	"golang.org/x/oauth2"
)

// TokenManager manages OAuth2 tokens
type TokenManager struct {
	tokens map[string]*oauth2.Token
	mu     sync.RWMutex
}

// TokenEntry represents a stored token with metadata
type TokenEntry struct {
	Token     *oauth2.Token
	UserID    string
	Provider  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewTokenManager creates a new token manager
func NewTokenManager() *TokenManager {
	return &TokenManager{
		tokens: make(map[string]*oauth2.Token),
	}
}

// Store stores a token
func (tm *TokenManager) Store(provider, userID string, token *oauth2.Token) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	key := tm.makeKey(provider, userID)
	tm.tokens[key] = token
}

// Get retrieves a token
func (tm *TokenManager) Get(provider, userID string) *oauth2.Token {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	key := tm.makeKey(provider, userID)
	return tm.tokens[key]
}

// Delete removes a token
func (tm *TokenManager) Delete(provider, userID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	key := tm.makeKey(provider, userID)
	delete(tm.tokens, key)
}

// IsExpired checks if token is expired
func (tm *TokenManager) IsExpired(provider, userID string) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	key := tm.makeKey(provider, userID)
	token, exists := tm.tokens[key]
	if !exists {
		return true
	}

	return token.Expiry.Before(time.Now())
}

// RefreshIfNeeded refreshes token if expired
func (tm *TokenManager) RefreshIfNeeded(provider, userID string, newToken *oauth2.Token) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	key := tm.makeKey(provider, userID)
	tm.tokens[key] = newToken
}

// Clear clears all tokens
func (tm *TokenManager) Clear() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.tokens = make(map[string]*oauth2.Token)
}

// Count returns number of stored tokens
func (tm *TokenManager) Count() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	return len(tm.tokens)
}

func (tm *TokenManager) makeKey(provider, userID string) string {
	return provider + ":" + userID
}
