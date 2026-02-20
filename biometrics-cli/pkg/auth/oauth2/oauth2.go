// Package oauth2 provides OAuth2 client implementation with support for
// multiple providers (Google, GitHub) and secure token management.
//
// This package implements:
//   - OAuth2 client flow
//   - Token management and refresh
//   - Provider integration (Google, GitHub)
//   - Callback handlers
//   - Session management
//   - Security best practices
//
// Best Practices Feb 2026 compliant.
package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// OAuth2 errors
var (
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidToken     = errors.New("invalid token")
	ErrProviderNotFound = errors.New("provider not found")
	ErrStateMismatch    = errors.New("state mismatch")
	ErrCodeRequired     = errors.New("authorization code required")
)

// ProviderType represents supported OAuth2 providers
type ProviderType string

const (
	ProviderGoogle ProviderType = "google"
	ProviderGitHub ProviderType = "github"
)

// OAuth2Client manages OAuth2 authentication
type OAuth2Client struct {
	configs     map[ProviderType]*oauth2.Config
	tokens      *TokenManager
	httpClient  *http.Client
	stateSecret string
	mu          sync.RWMutex
	redirectURL string
	scopes      map[ProviderType][]string
}

// OAuth2Config holds OAuth2 configuration
type OAuth2Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Provider     ProviderType
}

// UserInfo contains user information from OAuth2 provider
type UserInfo struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Provider    string `json:"provider"`
	AccessToken string `json:"access_token"`
}

// NewOAuth2Client creates a new OAuth2 client
func NewOAuth2Client(configs []OAuth2Config) (*OAuth2Client, error) {
	client := &OAuth2Client{
		configs:     make(map[ProviderType]*oauth2.Config),
		tokens:      NewTokenManager(),
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		stateSecret: generateStateSecret(),
		scopes:      make(map[ProviderType][]string),
	}

	for _, cfg := range configs {
		if err := client.AddProvider(cfg); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// AddProvider adds an OAuth2 provider
func (c *OAuth2Client) AddProvider(cfg OAuth2Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var oauthConfig *oauth2.Config

	switch cfg.Provider {
	case ProviderGoogle:
		oauthConfig = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     google.Endpoint,
		}
	case ProviderGitHub:
		oauthConfig = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     github.Endpoint,
		}
	default:
		return fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	c.configs[cfg.Provider] = oauthConfig
	c.scopes[cfg.Provider] = cfg.Scopes
	c.redirectURL = cfg.RedirectURL

	return nil
}

// GetAuthURL generates OAuth2 authorization URL
func (c *OAuth2Client) GetAuthURL(provider ProviderType) (string, string, error) {
	c.mu.RLock()
	config, exists := c.configs[provider]
	c.mu.RUnlock()

	if !exists {
		return "", "", ErrProviderNotFound
	}

	state := generateState()
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return authURL, state, nil
}

// Exchange exchanges authorization code for token
func (c *OAuth2Client) Exchange(ctx context.Context, provider ProviderType, code string) (*oauth2.Token, error) {
	c.mu.RLock()
	config, exists := c.configs[provider]
	c.mu.RUnlock()

	if !exists {
		return nil, ErrProviderNotFound
	}

	if code == "" {
		return nil, ErrCodeRequired
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, nil
}

// GetUserinfo retrieves user information from provider
func (c *OAuth2Client) GetUserInfo(ctx context.Context, provider ProviderType, token *oauth2.Token) (*UserInfo, error) {
	c.mu.RLock()
	config, exists := c.configs[provider]
	c.mu.RUnlock()

	if !exists {
		return nil, ErrProviderNotFound
	}

	client := config.Client(ctx, token)

	var userInfoURL string
	switch provider {
	case ProviderGoogle:
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	case ProviderGitHub:
		userInfoURL = "https://api.github.com/user"
	default:
		return nil, ErrProviderNotFound
	}

	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	userInfo.Provider = string(provider)
	userInfo.AccessToken = token.AccessToken

	return &userInfo, nil
}

// RefreshToken refreshes an expired token
func (c *OAuth2Client) RefreshToken(ctx context.Context, provider ProviderType, refreshToken string) (*oauth2.Token, error) {
	c.mu.RLock()
	config, exists := c.configs[provider]
	c.mu.RUnlock()

	if !exists {
		return nil, ErrProviderNotFound
	}

	tokenSource := config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

// RevokeToken revokes an OAuth2 token
func (c *OAuth2Client) RevokeToken(ctx context.Context, provider ProviderType, token string) error {
	c.mu.RLock()
	_, exists := c.configs[provider]
	c.mu.RUnlock()

	if !exists {
		return ErrProviderNotFound
	}

	var revokeURL string
	switch provider {
	case ProviderGoogle:
		revokeURL = fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", token)
	case ProviderGitHub:
		return nil
	default:
		return ErrProviderNotFound
	}

	resp, err := c.httpClient.Post(revokeURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// ValidateToken validates an OAuth2 token
func (c *OAuth2Client) ValidateToken(ctx context.Context, provider ProviderType, token string) (bool, error) {
	c.mu.RLock()
	_, exists := c.configs[provider]
	c.mu.RUnlock()

	if !exists {
		return false, ErrProviderNotFound
	}

	var validateURL string
	switch provider {
	case ProviderGoogle:
		validateURL = fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?access_token=%s", token)
	case ProviderGitHub:
		validateURL = "https://api.github.com/user"
	default:
		return false, ErrProviderNotFound
	}

	req, err := http.NewRequestWithContext(ctx, "GET", validateURL, nil)
	if err != nil {
		return false, err
	}

	if provider == ProviderGitHub {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// GetTokenManager returns the token manager
func (c *OAuth2Client) GetTokenManager() *TokenManager {
	return c.tokens
}

// generateState generates a secure state parameter
func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// generateStateSecret generates state secret for validation
func generateStateSecret() string {
	b := make([]byte, 64)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// ValidateState validates OAuth2 state parameter
func (c *OAuth2Client) ValidateState(state, expected string) error {
	if state == "" || expected == "" {
		return ErrStateMismatch
	}

	if !strings.EqualFold(state, expected) {
		return ErrStateMismatch
	}

	return nil
}

// GetProviderConfig returns OAuth2 config for provider
func (c *OAuth2Client) GetProviderConfig(provider ProviderType) (*oauth2.Config, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cfg, exists := c.configs[provider]; exists {
		return cfg, nil
	}

	return nil, ErrProviderNotFound
}

// ListProviders returns list of configured providers
func (c *OAuth2Client) ListProviders() []ProviderType {
	c.mu.RLock()
	defer c.mu.RUnlock()

	providers := make([]ProviderType, 0, len(c.configs))
	for provider := range c.configs {
		providers = append(providers, provider)
	}

	return providers
}

// ParseCallback parses OAuth2 callback URL
func (c *OAuth2Client) ParseCallback(callbackURL string) (string, string, error) {
	u, err := url.Parse(callbackURL)
	if err != nil {
		return "", "", err
	}

	query := u.Query()
	code := query.Get("code")
	state := query.Get("state")

	if code == "" {
		if errMsg := query.Get("error"); errMsg != "" {
			return "", "", fmt.Errorf("oauth2 error: %s", query.Get("error_description"))
		}
		return "", "", ErrCodeRequired
	}

	return code, state, nil
}

// CompleteAuth completes OAuth2 authentication flow
func (c *OAuth2Client) CompleteAuth(ctx context.Context, provider ProviderType, callbackURL, expectedState string) (*UserInfo, error) {
	code, state, err := c.ParseCallback(callbackURL)
	if err != nil {
		return nil, err
	}

	if err := c.ValidateState(state, expectedState); err != nil {
		return nil, err
	}

	token, err := c.Exchange(ctx, provider, code)
	if err != nil {
		return nil, err
	}

	userInfo, err := c.GetUserInfo(ctx, provider, token)
	if err != nil {
		return nil, err
	}

	c.tokens.Store(string(provider), userInfo.ID, token)

	return userInfo, nil
}

type TokenRefresher struct {
	client        *OAuth2Client
	provider      ProviderType
	userID        string
	refreshToken  string
	tokenMu       sync.RWMutex
	currentToken  *oauth2.Token
	stopChan      chan struct{}
	refreshBefore time.Duration
	interval      time.Duration
}

type TokenRefresherOption func(*TokenRefresher)

func WithRefreshInterval(interval time.Duration) TokenRefresherOption {
	return func(tr *TokenRefresher) {
		tr.interval = interval
	}
}

func WithRefreshBefore(duration time.Duration) TokenRefresherOption {
	return func(tr *TokenRefresher) {
		tr.refreshBefore = duration
	}
}

func (c *OAuth2Client) NewTokenRefresher(provider ProviderType, userID, refreshToken string, opts ...TokenRefresherOption) (*TokenRefresher, error) {
	tr := &TokenRefresher{
		client:        c,
		provider:      provider,
		userID:        userID,
		refreshToken:  refreshToken,
		stopChan:      make(chan struct{}),
		refreshBefore: 5 * time.Minute,
		interval:      1 * time.Minute,
	}

	for _, opt := range opts {
		opt(tr)
	}

	return tr, nil
}

func (tr *TokenRefresher) Start(ctx context.Context) error {
	go tr.run(ctx)
	return nil
}

func (tr *TokenRefresher) Stop() {
	close(tr.stopChan)
}

func (tr *TokenRefresher) GetToken(ctx context.Context) (*oauth2.Token, error) {
	tr.tokenMu.RLock()
	token := tr.currentToken
	tr.tokenMu.RUnlock()

	if token == nil {
		return nil, ErrInvalidToken
	}

	if tr.needsRefresh(token) {
		return tr.Refresh(ctx)
	}

	return token, nil
}

func (tr *TokenRefresher) Refresh(ctx context.Context) (*oauth2.Token, error) {
	tr.tokenMu.Lock()
	defer tr.tokenMu.Unlock()

	newToken, err := tr.client.RefreshToken(ctx, tr.provider, tr.refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	if newToken.RefreshToken != "" {
		tr.refreshToken = newToken.RefreshToken
	}

	tr.currentToken = newToken
	tr.client.tokens.Store(string(tr.provider), tr.userID, newToken)

	return newToken, nil
}

func (tr *TokenRefresher) needsRefresh(token *oauth2.Token) bool {
	if token == nil {
		return true
	}

	expiry := token.Expiry
	if expiry.IsZero() {
		return false
	}

	return time.Now().Add(tr.refreshBefore).After(expiry)
}

func (tr *TokenRefresher) run(ctx context.Context) {
	ticker := time.NewTicker(tr.interval)
	defer ticker.Stop()

	for {
		select {
		case <-tr.stopChan:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			tr.tokenMu.RLock()
			token := tr.currentToken
			tr.tokenMu.RUnlock()

			if token != nil && tr.needsRefresh(token) {
				newToken, err := tr.Refresh(ctx)
				if err != nil {
					continue
				}
				tr.client.tokens.Store(string(tr.provider), tr.userID, newToken)
			}
		}
	}
}

func (tr *TokenRefresher) IsTokenExpiredOrExpiring() bool {
	tr.tokenMu.RLock()
	defer tr.tokenMu.RUnlock()

	if tr.currentToken == nil {
		return true
	}

	return tr.needsRefresh(tr.currentToken)
}

func (tr *TokenRefresher) GetRefreshToken() string {
	tr.tokenMu.RLock()
	defer tr.tokenMu.RUnlock()
	return tr.refreshToken
}

func (tr *TokenRefresher) UpdateRefreshToken(newRefreshToken string) {
	tr.tokenMu.Lock()
	defer tr.tokenMu.Unlock()
	tr.refreshToken = newRefreshToken
}

func (c *OAuth2Client) Logout(ctx context.Context, provider ProviderType, userID string) error {
	token := c.tokens.Get(string(provider), userID)
	if token != nil {
		if err := c.RevokeToken(ctx, provider, token.AccessToken); err != nil {
			return err
		}
	}

	c.tokens.Delete(string(provider), userID)
	return nil
}

// GetHTTPClient returns configured HTTP client
func (c *OAuth2Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// SetHTTPClient sets HTTP client
func (c *OAuth2Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// IsConfigured checks if provider is configured
func (c *OAuth2Client) IsConfigured(provider ProviderType) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.configs[provider]
	return exists
}

// GetRedirectURL returns configured redirect URL
func (c *OAuth2Client) GetRedirectURL() string {
	return c.redirectURL
}
