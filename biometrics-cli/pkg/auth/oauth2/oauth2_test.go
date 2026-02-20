package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

// TestOAuth2ClientCreation tests OAuth2 client initialization
func TestOAuth2ClientCreation(t *testing.T) {
	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"openid", "email", "profile"},
			Provider:     ProviderGoogle,
		},
		{
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://localhost:8080/callback",
			Scopes:       []string{"user:email", "read:user"},
			Provider:     ProviderGitHub,
		},
	}

	client, err := NewOAuth2Client(configs)
	if err != nil {
		t.Fatalf("Failed to create OAuth2 client: %v", err)
	}

	if client == nil {
		t.Fatal("Expected OAuth2 client, got nil")
	}

	if len(client.configs) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(client.configs))
	}

	if client.httpClient == nil {
		t.Error("Expected HTTP client to be initialized")
	}

	if client.tokens == nil {
		t.Error("Expected token manager to be initialized")
	}
}

// TestAddProvider tests adding OAuth2 providers
func TestAddProvider(t *testing.T) {
	client := &OAuth2Client{
		configs:    make(map[ProviderType]*oauth2.Config),
		tokens:     NewTokenManager(),
		httpClient: &http.Client{Timeout: 30 * time.Second},
		scopes:     make(map[ProviderType][]string),
	}

	tests := []struct {
		name        string
		provider    ProviderType
		expectError bool
	}{
		{"Google Provider", ProviderGoogle, false},
		{"GitHub Provider", ProviderGitHub, false},
		{"Invalid Provider", ProviderType("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := OAuth2Config{
				ClientID:     "test-id",
				ClientSecret: "test-secret",
				RedirectURL:  "http://localhost/callback",
				Scopes:       []string{"openid"},
				Provider:     tt.provider,
			}

			err := client.AddProvider(cfg)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error=%v, got %v", tt.expectError, err)
			}
		})
	}
}

// TestGetAuthURL tests authorization URL generation
func TestGetAuthURL(t *testing.T) {
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

	authURL, state, err := client.GetAuthURL(ProviderGoogle)
	if err != nil {
		t.Fatalf("GetAuthURL failed: %v", err)
	}

	if authURL == "" {
		t.Error("Expected non-empty auth URL")
	}

	if state == "" {
		t.Error("Expected non-empty state parameter")
	}

	_, _, err = client.GetAuthURL(ProviderType("invalid"))
	if err != ErrProviderNotFound {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
}

// TestExchange tests authorization code exchange
func TestExchange(t *testing.T) {
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

	_, err = client.Exchange(context.Background(), ProviderGoogle, "")
	if err != ErrCodeRequired {
		t.Errorf("Expected ErrCodeRequired, got %v", err)
	}

	_, err = client.Exchange(context.Background(), ProviderType("invalid"), "test-code")
	if err != ErrProviderNotFound {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
}

// TestGetUserInfo tests user info retrieval
func TestGetUserInfo(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := UserInfo{
			ID:    "123456",
			Email: "test@example.com",
			Name:  "Test User",
		}
		json.NewEncoder(w).Encode(userInfo)
	}))
	defer mockServer.Close()

	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  mockServer.URL,
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, err := NewOAuth2Client(configs)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Override HTTP client to use mock server
	client.httpClient = mockServer.Client()

	token := &oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
	}

	userInfo, err := client.GetUserInfo(context.Background(), ProviderGoogle, token)
	if err != nil {
		t.Fatalf("GetUserInfo failed: %v", err)
	}

	if userInfo.ID != "123456" {
		t.Errorf("Expected ID 123456, got %s", userInfo.ID)
	}
}

// TestRefreshToken tests token refresh
func TestRefreshToken(t *testing.T) {
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

	_, err = client.RefreshToken(context.Background(), ProviderType("invalid"), "refresh-token")
	if err != ErrProviderNotFound {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
}

// TestRevokeToken tests token revocation
func TestRevokeToken(t *testing.T) {
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

	err = client.RevokeToken(context.Background(), ProviderType("invalid"), "token")
	if err != ErrProviderNotFound {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
}

// TestValidateToken tests token validation
func TestValidateToken(t *testing.T) {
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

	valid, err := client.ValidateToken(context.Background(), ProviderType("invalid"), "token")
	if err != ErrProviderNotFound {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
	if valid {
		t.Error("Expected token to be invalid")
	}
}

// TestValidateState tests state parameter validation
func TestValidateState(t *testing.T) {
	client := &OAuth2Client{}

	tests := []struct {
		name     string
		state    string
		expected string
		wantErr  bool
	}{
		{"Valid state", "abc123", "abc123", false},
		{"Mismatch", "abc123", "xyz789", true},
		{"Empty state", "", "abc123", true},
		{"Empty expected", "abc123", "", true},
		{"Case insensitive", "ABC123", "abc123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ValidateState(tt.state, tt.expected)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, got %v", tt.wantErr, err)
			}
			if tt.wantErr && err != ErrStateMismatch {
				t.Errorf("Expected ErrStateMismatch, got %v", err)
			}
		})
	}
}

// TestParseCallback tests callback URL parsing
func TestParseCallback(t *testing.T) {
	client := &OAuth2Client{}

	tests := []struct {
		name      string
		callback  string
		wantCode  string
		wantState string
		wantErr   bool
	}{
		{
			name:      "Valid callback",
			callback:  "http://localhost/callback?code=abc123&state=xyz789",
			wantCode:  "abc123",
			wantState: "xyz789",
			wantErr:   false,
		},
		{
			name:     "Missing code",
			callback: "http://localhost/callback?state=xyz789",
			wantErr:  true,
		},
		{
			name:     "Error in callback",
			callback: "http://localhost/callback?error=access_denied&error_description=User+denied+access",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, state, err := client.ParseCallback(tt.callback)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr {
				if code != tt.wantCode {
					t.Errorf("Expected code %s, got %s", tt.wantCode, code)
				}
				if state != tt.wantState {
					t.Errorf("Expected state %s, got %s", tt.wantState, state)
				}
			}
		})
	}
}

// TestListProviders tests listing configured providers
func TestListProviders(t *testing.T) {
	client := &OAuth2Client{
		configs: map[ProviderType]*oauth2.Config{
			ProviderGoogle: {},
			ProviderGitHub: {},
		},
	}

	providers := client.ListProviders()
	if len(providers) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(providers))
	}
}

// TestIsConfigured tests provider configuration check
func TestIsConfigured(t *testing.T) {
	client := &OAuth2Client{
		configs: map[ProviderType]*oauth2.Config{
			ProviderGoogle: {},
		},
	}

	if !client.IsConfigured(ProviderGoogle) {
		t.Error("Expected Google to be configured")
	}

	if client.IsConfigured(ProviderGitHub) {
		t.Error("Expected GitHub to not be configured")
	}
}

// TestGetRedirectURL tests redirect URL retrieval
func TestGetRedirectURL(t *testing.T) {
	client := &OAuth2Client{
		redirectURL: "http://localhost:8080/callback",
	}

	url := client.GetRedirectURL()
	if url != "http://localhost:8080/callback" {
		t.Errorf("Expected redirect URL http://localhost:8080/callback, got %s", url)
	}
}

// TestTokenManager tests token manager operations
func TestTokenManager(t *testing.T) {
	tm := NewTokenManager()

	token := &oauth2.Token{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	tm.Store("google", "user123", token)

	retrieved := tm.Get("google", "user123")
	if retrieved == nil {
		t.Fatal("Expected to retrieve token")
	}
	if retrieved.AccessToken != "access-token" {
		t.Errorf("Expected access-token, got %s", retrieved.AccessToken)
	}

	if tm.IsExpired("google", "user123") {
		t.Error("Expected token to not be expired")
	}

	if tm.Count() != 1 {
		t.Errorf("Expected 1 token, got %d", tm.Count())
	}

	tm.Delete("google", "user123")
	if tm.Get("google", "user123") != nil {
		t.Error("Expected token to be deleted")
	}

	tm.Store("google", "user1", &oauth2.Token{})
	tm.Store("google", "user2", &oauth2.Token{})
	tm.Clear()
	if tm.Count() != 0 {
		t.Errorf("Expected 0 tokens after clear, got %d", tm.Count())
	}
}

// TestRefreshIfNeeded tests token refresh logic
func TestRefreshIfNeeded(t *testing.T) {
	tm := NewTokenManager()

	oldToken := &oauth2.Token{
		AccessToken:  "old-token",
		RefreshToken: "old-refresh",
	}
	tm.Store("google", "user123", oldToken)

	newToken := &oauth2.Token{
		AccessToken:  "new-token",
		RefreshToken: "new-refresh",
	}
	tm.RefreshIfNeeded("google", "user123", newToken)

	retrieved := tm.Get("google", "user123")
	if retrieved.AccessToken != "new-token" {
		t.Errorf("Expected new-token, got %s", retrieved.AccessToken)
	}
}

// TestGetHTTPClient tests HTTP client getter/setter
func TestGetHTTPClient(t *testing.T) {
	client := &OAuth2Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	httpClient := client.GetHTTPClient()
	if httpClient == nil {
		t.Error("Expected HTTP client")
	}

	newClient := &http.Client{Timeout: 60 * time.Second}
	client.SetHTTPClient(newClient)
	if client.httpClient.Timeout != 60*time.Second {
		t.Error("Expected HTTP client to be updated")
	}
}

// TestGenerateState tests state generation
func TestGenerateState(t *testing.T) {
	state1 := generateState()
	state2 := generateState()

	if state1 == state2 {
		t.Error("Expected different state values")
	}

	if len(state1) == 0 {
		t.Error("Expected non-empty state")
	}
}

// TestGenerateStateSecret tests state secret generation
func TestGenerateStateSecret(t *testing.T) {
	secret1 := generateStateSecret()
	secret2 := generateStateSecret()

	if secret1 == secret2 {
		t.Error("Expected different secrets")
	}

	if len(secret1) == 0 {
		t.Error("Expected non-empty secret")
	}
}

// TestGetProviderConfig tests provider config retrieval
func TestGetProviderConfig(t *testing.T) {
	client := &OAuth2Client{
		configs: map[ProviderType]*oauth2.Config{
			ProviderGoogle: {ClientID: "google-id"},
		},
	}

	cfg, err := client.GetProviderConfig(ProviderGoogle)
	if err != nil {
		t.Fatalf("GetProviderConfig failed: %v", err)
	}
	if cfg.ClientID != "google-id" {
		t.Errorf("Expected google-id, got %s", cfg.ClientID)
	}

	_, err = client.GetProviderConfig(ProviderGitHub)
	if err != ErrProviderNotFound {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
}

// TestLogout tests logout functionality
func TestLogout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  mockServer.URL,
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, err := NewOAuth2Client(configs)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	token := &oauth2.Token{AccessToken: "test-token"}
	client.tokens.Store("google", "user123", token)

	err = client.Logout(context.Background(), ProviderGoogle, "user123")
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	if client.tokens.Get("google", "user123") != nil {
		t.Error("Expected token to be deleted after logout")
	}
}

// TestOAuth2FlowIntegration tests complete OAuth2 flow with mock server
func TestOAuth2FlowIntegration(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"access_token":"mock-token","token_type":"Bearer","expires_in":3600}`)
		case "/userinfo":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"id":"123","email":"test@example.com","name":"Test User"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	configs := []OAuth2Config{
		{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  mockServer.URL + "/callback",
			Scopes:       []string{"openid"},
			Provider:     ProviderGoogle,
		},
	}

	client, err := NewOAuth2Client(configs)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	authURL, state, err := client.GetAuthURL(ProviderGoogle)
	if err != nil {
		t.Fatalf("GetAuthURL failed: %v", err)
	}

	if authURL == "" {
		t.Error("Expected non-empty auth URL")
	}

	callbackURL := fmt.Sprintf("%s/callback?code=test-auth-code&state=%s", mockServer.URL, state)
	code, parsedState, err := client.ParseCallback(callbackURL)
	if err != nil {
		t.Fatalf("ParseCallback failed: %v", err)
	}

	if err := client.ValidateState(parsedState, state); err != nil {
		t.Fatalf("ValidateState failed: %v", err)
	}

	if code != "test-auth-code" {
		t.Errorf("Expected code test-auth-code, got %s", code)
	}
}

// TestConcurrentAccess tests thread safety
func TestConcurrentAccess(t *testing.T) {
	client := &OAuth2Client{
		configs:    make(map[ProviderType]*oauth2.Config),
		tokens:     NewTokenManager(),
		httpClient: &http.Client{},
		scopes:     make(map[ProviderType][]string),
	}

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			client.ListProviders()
			client.IsConfigured(ProviderGoogle)
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
