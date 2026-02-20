package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// Provider defines interface for OAuth2 providers
type Provider interface {
	Name() ProviderType
	Endpoint() oauth2.Endpoint
	Scopes() []string
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error)
}

// GoogleProvider implements Google OAuth2
type GoogleProvider struct {
	config *oauth2.Config
}

// GitHubProvider implements GitHub OAuth2
type GitHubProvider struct {
	config *oauth2.Config
}

// NewGoogleProvider creates Google provider
func NewGoogleProvider(clientID, clientSecret, redirectURL string, scopes []string) *GoogleProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     google.Endpoint,
		},
	}
}

// NewGitHubProvider creates GitHub provider
func NewGitHubProvider(clientID, clientSecret, redirectURL string, scopes []string) *GitHubProvider {
	return &GitHubProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     github.Endpoint,
		},
	}
}

// Name returns provider name
func (p *GoogleProvider) Name() ProviderType {
	return ProviderGoogle
}

// Name returns provider name
func (p *GitHubProvider) Name() ProviderType {
	return ProviderGitHub
}

// Endpoint returns OAuth2 endpoint
func (p *GoogleProvider) Endpoint() oauth2.Endpoint {
	return p.config.Endpoint
}

// Endpoint returns OAuth2 endpoint
func (p *GitHubProvider) Endpoint() oauth2.Endpoint {
	return p.config.Endpoint
}

// Scopes returns configured scopes
func (p *GoogleProvider) Scopes() []string {
	return p.config.Scopes
}

// Scopes returns configured scopes
func (p *GitHubProvider) Scopes() []string {
	return p.config.Scopes
}

func parseUserInfo(body io.Reader, provider, accessToken string) (*UserInfo, error) {
	var userInfo UserInfo
	if err := json.NewDecoder(body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}
	userInfo.Provider = provider
	userInfo.AccessToken = accessToken
	return &userInfo, nil
}

// GetAuthURL generates authorization URL
func (p *GoogleProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// GetAuthURL generates authorization URL
func (p *GitHubProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

// Exchange exchanges code for token
func (p *GoogleProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

// Exchange exchanges code for token
func (p *GitHubProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.config.Exchange(ctx, code)
}

// RefreshToken refreshes token
func (p *GoogleProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	ts := p.config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	return ts.Token()
}

// RefreshToken refreshes token
func (p *GitHubProvider) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	ts := p.config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	return ts.Token()
}
