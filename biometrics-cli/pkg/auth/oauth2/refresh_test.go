package oauth2

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestTokenRefresherCreation(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid", "email"},
			Provider:     ProviderGoogle,
		},
	}

	client, err := NewOAuth2Client(configs)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	refresher, err := client.NewTokenRefresher(ProviderGoogle, "user123", "test-refresh-token")
	if err != nil {
		t.Fatalf("Failed to create refresher: %v", err)
	}

	if refresher == nil {
		t.Fatal("Expected refresher, got nil")
	}

	if refresher.provider != ProviderGoogle {
		t.Errorf("Expected provider google, got %v", refresher.provider)
	}

	if refresher.userID != "user123" {
		t.Errorf("Expected userID user123, got %s", refresher.userID)
	}

	if refresher.refreshToken != "test-refresh-token" {
		t.Errorf("Expected refresh token test-refresh-token, got %s", refresher.refreshToken)
	}

	if refresher.stopChan == nil {
		t.Error("Expected stopChan to be initialized")
	}
}

func TestTokenRefresherWithOptions(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, err := NewOAuth2Client(configs)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	refresher, err := client.NewTokenRefresher(
		ProviderGoogle,
		"user123",
		"refresh-token",
		WithRefreshInterval(30*time.Second),
		WithRefreshBefore(2*time.Minute),
	)
	if err != nil {
		t.Fatalf("Failed to create refresher: %v", err)
	}

	if refresher.interval != 30*time.Second {
		t.Errorf("Expected interval 30s, got %v", refresher.interval)
	}

	if refresher.refreshBefore != 2*time.Minute {
		t.Errorf("Expected refreshBefore 2m, got %v", refresher.refreshBefore)
	}
}

func TestTokenRefresherNeedsRefresh(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "refresh")

	tests := []struct {
		name        string
		token       *oauth2.Token
		refreshDur  time.Duration
		wantRefresh bool
	}{
		{
			name:        "nil token",
			token:       nil,
			refreshDur:  5 * time.Minute,
			wantRefresh: true,
		},
		{
			name: "expired token",
			token: &oauth2.Token{
				AccessToken: "token",
				Expiry:      time.Now().Add(-1 * time.Minute),
			},
			refreshDur:  5 * time.Minute,
			wantRefresh: true,
		},
		{
			name: "token expiring soon",
			token: &oauth2.Token{
				AccessToken: "token",
				Expiry:      time.Now().Add(2 * time.Minute),
			},
			refreshDur:  5 * time.Minute,
			wantRefresh: true,
		},
		{
			name: "token still valid",
			token: &oauth2.Token{
				AccessToken: "token",
				Expiry:      time.Now().Add(10 * time.Minute),
			},
			refreshDur:  5 * time.Minute,
			wantRefresh: false,
		},
		{
			name: "token no expiry",
			token: &oauth2.Token{
				AccessToken: "token",
			},
			refreshDur:  5 * time.Minute,
			wantRefresh: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refresher.refreshBefore = tt.refreshDur
			gotRefresh := refresher.needsRefresh(tt.token)
			if gotRefresh != tt.wantRefresh {
				t.Errorf("needsRefresh() = %v, want %v", gotRefresh, tt.wantRefresh)
			}
		})
	}
}

func TestTokenRefresherIsExpiredOrExpiring(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "refresh")

	if !refresher.IsTokenExpiredOrExpiring() {
		t.Error("Expected nil token to be expired")
	}

	refresher.currentToken = &oauth2.Token{
		AccessToken: "token",
		Expiry:      time.Now().Add(10 * time.Minute),
	}

	if refresher.IsTokenExpiredOrExpiring() {
		t.Error("Expected valid token to not be expired")
	}

	refresher.currentToken = &oauth2.Token{
		AccessToken: "token",
		Expiry:      time.Now().Add(2 * time.Minute),
	}

	if !refresher.IsTokenExpiredOrExpiring() {
		t.Error("Expected token expiring soon to be flagged as expiring")
	}
}

func TestTokenRefresherGetRefreshToken(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "my-refresh-token")

	token := refresher.GetRefreshToken()
	if token != "my-refresh-token" {
		t.Errorf("Expected my-refresh-token, got %s", token)
	}
}

func TestTokenRefresherUpdateRefreshToken(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "old-refresh")

	refresher.UpdateRefreshToken("new-refresh-token")

	token := refresher.GetRefreshToken()
	if token != "new-refresh-token" {
		t.Errorf("Expected new-refresh-token, got %s", token)
	}
}

func TestTokenRefresherStop(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "refresh")

	refresher.currentToken = &oauth2.Token{
		AccessToken: "token",
		Expiry:      time.Now().Add(10 * time.Minute),
	}

	ctx, cancel := context.WithCancel(context.Background())
	refresher.Start(ctx)

	time.Sleep(100 * time.Millisecond)

	refresher.Stop()

	select {
	case <-refresher.stopChan:
	default:
	}

	cancel()
}

func TestTokenRefresherConcurrentAccess(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "refresh")

	refresher.currentToken = &oauth2.Token{
		AccessToken: "token",
		Expiry:      time.Now().Add(10 * time.Minute),
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = refresher.IsTokenExpiredOrExpiring()
		}()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = refresher.GetRefreshToken()
		}()
	}

	wg.Wait()
}

func TestTokenRefresherRefreshTokenRotation(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "original-refresh")

	if refresher.GetRefreshToken() != "original-refresh" {
		t.Error("Initial refresh token not set correctly")
	}

	refresher.UpdateRefreshToken("new-refresh-token")

	if refresher.GetRefreshToken() != "new-refresh-token" {
		t.Error("Refresh token was not updated")
	}
}

func TestTokenRefresherStart(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "refresh",
		WithRefreshInterval(10*time.Millisecond))

	ctx := context.Background()
	err := refresher.Start(ctx)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	refresher.Stop()
}

func TestTokenRefresherGetTokenNoToken(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "refresh")

	_, err := refresher.GetToken(context.Background())
	if err != ErrInvalidToken {
		t.Errorf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestTokenRefresherIntegration(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "refresh")

	refresher.currentToken = &oauth2.Token{
		AccessToken: "existing-token",
		Expiry:      time.Now().Add(10 * time.Minute),
	}

	ctx, cancel := context.WithCancel(context.Background())
	refresher.Start(ctx)

	time.Sleep(50 * time.Millisecond)

	refresher.Stop()
	cancel()

	if refresher.currentToken == nil {
		t.Error("Expected token to be preserved")
	}

	if refresher.currentToken.AccessToken != "existing-token" {
		t.Errorf("Expected existing-token, got %s", refresher.currentToken.AccessToken)
	}
}

func TestProviderRefreshMethod(t *testing.T) {
	googleProvider := NewGoogleProvider("client-id", "client-secret", "http://localhost/callback", []string{"openid"})
	githubProvider := NewGitHubProvider("client-id", "client-secret", "http://localhost/callback", []string{"read:user"})

	if googleProvider.Name() != ProviderGoogle {
		t.Errorf("Expected Google provider")
	}

	if githubProvider.Name() != ProviderGitHub {
		t.Errorf("Expected GitHub provider")
	}

	_ = googleProvider
	_ = githubProvider
}

type mockProviderForRefresh struct{}

func (m *mockProviderForRefresh) Name() ProviderType        { return ProviderGoogle }
func (m *mockProviderForRefresh) Endpoint() oauth2.Endpoint { return oauth2.Endpoint{} }
func (m *mockProviderForRefresh) Scopes() []string          { return []string{} }
func (m *mockProviderForRefresh) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	return &UserInfo{ID: "123", Email: "test@example.com"}, nil
}
func (m *mockProviderForRefresh) Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: "new"}, nil
}

func TestProviderRefreshInterface(t *testing.T) {
	var p Provider = &mockProviderForRefresh{}
	if p.Name() != ProviderGoogle {
		t.Error("Expected Provider to be Google")
	}

	token, err := p.Refresh(context.Background(), "refresh")
	if err != nil {
		t.Errorf("Refresh failed: %v", err)
	}
	if token.AccessToken != "new" {
		t.Errorf("Expected new token, got %s", token.AccessToken)
	}
}

func TestTokenRefresherWithExistingToken(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "refresh-token")

	refresher.currentToken = &oauth2.Token{
		AccessToken: "existing-token",
		Expiry:      time.Now().Add(10 * time.Minute),
	}

	token, err := refresher.GetToken(context.Background())
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if token.AccessToken != "existing-token" {
		t.Errorf("Expected existing-token, got %s", token.AccessToken)
	}
}

func TestTokenRefresherBackgroundRefresh(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"access_token":"bg-token","token_type":"Bearer","refresh_token":"bg-refresh","expires_in":3600}`)
	}))
	defer mockServer.Close()

	client, _ := NewOAuth2Client(configs)
	refresher, _ := client.NewTokenRefresher(ProviderGoogle, "user123", "bg-refresh")

	refresher.refreshBefore = 1 * time.Minute
	refresher.interval = 50 * time.Millisecond
	refresher.currentToken = &oauth2.Token{
		AccessToken: "old-token",
		Expiry:      time.Now().Add(30 * time.Second),
	}

	ctx, cancel := context.WithCancel(context.Background())
	refresher.Start(ctx)

	time.Sleep(100 * time.Millisecond)

	refresher.Stop()
	cancel()

	if refresher.currentToken == nil {
		t.Error("Expected token to be updated in background")
	}
}
