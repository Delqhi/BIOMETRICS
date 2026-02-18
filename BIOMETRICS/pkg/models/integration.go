package models

import (
	"time"

	"github.com/google/uuid"
)

type Integration struct {
	ID            uuid.UUID              `json:"id" db:"id"`
	TenantID      *uuid.UUID             `json:"tenant_id,omitempty" db:"tenant_id"`
	Name          string                 `json:"name" db:"name"`
	Slug          string                 `json:"slug" db:"slug"`
	Description   string                 `json:"description" db:"description"`
	Type          string                 `json:"type" db:"type"`
	Provider      string                 `json:"provider" db:"provider"`
	Status        string                 `json:"status" db:"status"`
	Logo          string                 `json:"logo" db:"logo"`
	Website       string                 `json:"website" db:"website"`
	OwnerID       uuid.UUID              `json:"owner_id" db:"owner_id"`
	APIKey        string                 `json:"-" db:"api_key"`
	APISecret     string                 `json:"-" db:"api_secret"`
	AccessToken   string                 `json:"-" db:"access_token"`
	RefreshToken  string                 `json:"-" db:"refresh_token"`
	WebhookURL    string                 `json:"webhook_url" db:"webhook_url"`
	WebhookSecret string                 `json:"-" db:"webhook_secret"`
	Scopes        []string               `json:"scopes" db:"scopes"`
	Settings      map[string]interface{} `json:"settings" db:"settings"`
	Credentials   map[string]interface{} `json:"credentials" db:"credentials"`
	RateLimit     int                    `json:"rate_limit" db:"rate_limit"`
	RateRemaining int                    `json:"rate_remaining" db:"rate_remaining"`
	RateResetAt   *time.Time             `json:"rate_reset_at" db:"rate_reset_at"`
	LastSyncAt    *time.Time             `json:"last_sync_at" db:"last_sync_at"`
	NextSyncAt    *time.Time             `json:"next_sync_at" db:"next_sync_at"`
	SyncInterval  int                    `json:"sync_interval" db:"sync_interval"`
	FailureCount  int                    `json:"failure_count" db:"failure_count"`
	LastError     string                 `json:"last_error" db:"last_error"`
	Metadata      map[string]interface{} `json:"metadata" db:"metadata"`
	IsActive      bool                   `json:"is_active" db:"is_active"`
	Version       int                    `json:"version" db:"version"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
}

type IntegrationType string

const (
	IntegrationTypeOAuth   IntegrationType = "oauth"
	IntegrationTypeAPI     IntegrationType = "api"
	IntegrationTypeWebhook IntegrationType = "webhook"
	IntegrationTypePlugin  IntegrationType = "plugin"
	IntegrationTypeCustom  IntegrationType = "custom"
)

type IntegrationProvider string

const (
	IntegrationProviderGoogle  IntegrationProvider = "google"
	IntegrationProviderGitHub  IntegrationProvider = "github"
	IntegrationProviderStripe  IntegrationProvider = "stripe"
	IntegrationProviderSlack   IntegrationProvider = "slack"
	IntegrationProviderDiscord IntegrationProvider = "discord"
	IntegrationProviderTwitter IntegrationProvider = "twitter"
	IntegrationProviderDropbox IntegrationProvider = "dropbox"
	IntegrationProviderCustom  IntegrationProvider = "custom"
)

type IntegrationStatus string

const (
	IntegrationStatusActive   IntegrationStatus = "active"
	IntegrationStatusInactive IntegrationStatus = "inactive"
	IntegrationStatusPending  IntegrationStatus = "pending"
	IntegrationStatusError    IntegrationStatus = "error"
	IntegrationStatusExpired  IntegrationStatus = "expired"
)

type CreateIntegrationInput struct {
	TenantID     string                 `json:"tenant_id"`
	Name         string                 `json:"name" binding:"required"`
	Slug         string                 `json:"slug"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type" binding:"required"`
	Provider     string                 `json:"provider" binding:"required"`
	Logo         string                 `json:"logo"`
	Website      string                 `json:"website"`
	WebhookURL   string                 `json:"webhook_url"`
	Scopes       []string               `json:"scopes"`
	Settings     map[string]interface{} `json:"settings"`
	SyncInterval int                    `json:"sync_interval"`
}

type UpdateIntegrationInput struct {
	Name         *string                `json:"name,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Status       *string                `json:"status,omitempty"`
	Logo         *string                `json:"logo,omitempty"`
	WebhookURL   *string                `json:"webhook_url,omitempty"`
	Scopes       []string               `json:"scopes,omitempty"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
	SyncInterval *int                   `json:"sync_interval,omitempty"`
	IsActive     *bool                  `json:"is_active,omitempty"`
}

type IntegrationFilter struct {
	Type     string `query:"type"`
	Provider string `query:"provider"`
	Status   string `query:"status"`
	TenantID string `query:"tenant_id"`
	OwnerID  string `query:"owner_id"`
	Search   string `query:"search"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

type IntegrationResponse struct {
	ID            uuid.UUID              `json:"id"`
	TenantID      *uuid.UUID             `json:"tenant_id,omitempty"`
	Name          string                 `json:"name"`
	Slug          string                 `json:"slug"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	Provider      string                 `json:"provider"`
	Status        string                 `json:"status"`
	Logo          string                 `json:"logo"`
	Website       string                 `json:"website"`
	OwnerID       uuid.UUID              `json:"owner_id"`
	Scopes        []string               `json:"scopes"`
	Settings      map[string]interface{} `json:"settings"`
	RateLimit     int                    `json:"rate_limit"`
	RateRemaining int                    `json:"rate_remaining"`
	RateResetAt   *time.Time             `json:"rate_reset_at"`
	LastSyncAt    *time.Time             `json:"last_sync_at"`
	NextSyncAt    *time.Time             `json:"next_sync_at"`
	FailureCount  int                    `json:"failure_count"`
	LastError     string                 `json:"last_error"`
	IsActive      bool                   `json:"is_active"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

func (i *Integration) ToResponse() *IntegrationResponse {
	return &IntegrationResponse{
		ID:            i.ID,
		TenantID:      i.TenantID,
		Name:          i.Name,
		Slug:          i.Slug,
		Description:   i.Description,
		Type:          i.Type,
		Provider:      i.Provider,
		Status:        i.Status,
		Logo:          i.Logo,
		Website:       i.Website,
		OwnerID:       i.OwnerID,
		Scopes:        i.Scopes,
		Settings:      i.Settings,
		RateLimit:     i.RateLimit,
		RateRemaining: i.RateRemaining,
		RateResetAt:   i.RateResetAt,
		LastSyncAt:    i.LastSyncAt,
		NextSyncAt:    i.NextSyncAt,
		FailureCount:  i.FailureCount,
		LastError:     i.LastError,
		IsActive:      i.IsActive,
		CreatedAt:     i.CreatedAt,
		UpdatedAt:     i.UpdatedAt,
	}
}
