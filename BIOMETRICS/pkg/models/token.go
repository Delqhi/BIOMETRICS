package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	UserID       uuid.UUID              `json:"user_id" db:"user_id"`
	TenantID     *uuid.UUID             `json:"tenant_id,omitempty" db:"tenant_id"`
	Type         string                 `json:"type" db:"type"`
	TokenHash    string                 `json:"-" db:"token_hash"`
	TokenPrefix  string                 `json:"token_prefix" db:"token_prefix"`
	RefreshToken string                 `json:"-" db:"refresh_token"`
	RefreshHash  string                 `json:"-" db:"refresh_hash"`
	IPAddress    string                 `json:"ip_address" db:"ip_address"`
	UserAgent    string                 `json:"user_agent" db:"user_agent"`
	DeviceID     string                 `json:"device_id" db:"device_id"`
	DeviceInfo   string                 `json:"device_info" db:"device_info"`
	Location     string                 `json:"location" db:"location"`
	Scopes       []string               `json:"scopes" db:"scopes"`
	Status       string                 `json:"status" db:"status"`
	IssuedAt     time.Time              `json:"issued_at" db:"issued_at"`
	ExpiresAt    time.Time              `json:"expires_at" db:"expires_at"`
	LastUsedAt   *time.Time             `json:"last_used_at" db:"last_used_at"`
	RevokedAt    *time.Time             `json:"revoked_at" db:"revoked_at"`
	RevokeReason string                 `json:"revoke_reason" db:"revoke_reason"`
	UseCount     int                    `json:"use_count" db:"use_count"`
	MaxUses      int                    `json:"max_uses" db:"max_uses"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	Version      int                    `json:"version" db:"version"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}

type TokenType string

const (
	TokenTypeAccess        TokenType = "access"
	TokenTypeRefresh       TokenType = "refresh"
	TokenTypeAPI           TokenType = "api"
	TokenTypeMagicLink     TokenType = "magic_link"
	TokenTypePasswordReset TokenType = "password_reset"
	TokenTypeEmailVerify   TokenType = "email_verify"
	TokenTypeInvite        TokenType = "invite"
	TokenTypeRecovery      TokenType = "recovery"
)

type TokenStatus string

const (
	TokenStatusActive    TokenStatus = "active"
	TokenStatusExpired   TokenStatus = "expired"
	TokenStatusRevoked   TokenStatus = "revoked"
	TokenStatusUsed      TokenStatus = "used"
	TokenStatusSuspended TokenStatus = "suspended"
)

type CreateTokenInput struct {
	UserID    string   `json:"user_id" binding:"required"`
	TenantID  string   `json:"tenant_id"`
	Type      string   `json:"type" binding:"required"`
	Scopes    []string `json:"scopes"`
	IPAddress string   `json:"ip_address"`
	UserAgent string   `json:"user_agent"`
	DeviceID  string   `json:"device_id"`
	ExpiresIn int      `json:"expires_in"`
	MaxUses   int      `json:"max_uses"`
}

type TokenFilter struct {
	UserID    string `query:"user_id"`
	TenantID  string `query:"tenant_id"`
	Type      string `query:"type"`
	Status    string `query:"status"`
	IPAddress string `query:"ip_address"`
	Page      int    `query:"page"`
	PageSize  int    `query:"page_size"`
}

type TokenResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty"`
	Type        string     `json:"type"`
	TokenPrefix string     `json:"token_prefix"`
	IPAddress   string     `json:"ip_address"`
	DeviceInfo  string     `json:"device_info"`
	Location    string     `json:"location"`
	Scopes      []string   `json:"scopes"`
	Status      string     `json:"status"`
	IssuedAt    time.Time  `json:"issued_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	UseCount    int        `json:"use_count"`
	MaxUses     int        `json:"max_uses"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (t *Token) ToResponse() *TokenResponse {
	return &TokenResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		TenantID:    t.TenantID,
		Type:        t.Type,
		TokenPrefix: t.TokenPrefix,
		IPAddress:   t.IPAddress,
		DeviceInfo:  t.DeviceInfo,
		Location:    t.Location,
		Scopes:      t.Scopes,
		Status:      t.Status,
		IssuedAt:    t.IssuedAt,
		ExpiresAt:   t.ExpiresAt,
		LastUsedAt:  t.LastUsedAt,
		UseCount:    t.UseCount,
		MaxUses:     t.MaxUses,
		CreatedAt:   t.CreatedAt,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
}
