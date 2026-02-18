package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID              `json:"id" db:"id"`
	Email               string                 `json:"email" db:"email"`
	PasswordHash        string                 `json:"-" db:"password_hash"`
	FirstName           string                 `json:"first_name" db:"first_name"`
	LastName            string                 `json:"last_name" db:"last_name"`
	Phone               string                 `json:"phone" db:"phone"`
	Role                string                 `json:"role" db:"role"`
	Status              string                 `json:"status" db:"status"`
	EmailVerified       bool                   `json:"email_verified" db:"email_verified"`
	TwoFactorEnabled    bool                   `json:"two_factor_enabled" db:"two_factor_enabled"`
	Provider            string                 `json:"provider" db:"provider"`
	ProviderID          string                 `json:"provider_id" db:"provider_id"`
	Avatar              string                 `json:"avatar" db:"avatar"`
	Bio                 string                 `json:"bio" db:"bio"`
	Roles               []string               `json:"roles" db:"roles"`
	Permissions         []string               `json:"permissions" db:"permissions"`
	TenantID            *uuid.UUID             `json:"tenant_id,omitempty" db:"tenant_id"`
	Metadata            map[string]interface{} `json:"metadata" db:"metadata"`
	LastLoginAt         *time.Time             `json:"last_login_at" db:"last_login_at"`
	FailedLoginAttempts int                    `json:"failed_login_attempts" db:"failed_login_attempts"`
	LockedUntil         *time.Time             `json:"locked_until" db:"locked_until"`
	CreatedAt           time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
}

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleUser     UserRole = "user"
	RoleGuest    UserRole = "guest"
	RoleOperator UserRole = "operator"
	RoleAuditor  UserRole = "auditor"
)

type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusPending   UserStatus = "pending"
	StatusSuspended UserStatus = "suspended"
	StatusLocked    UserStatus = "locked"
)

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:               u.ID,
		Email:            u.Email,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Phone:            u.Phone,
		Role:             u.Role,
		Status:           u.Status,
		EmailVerified:    u.EmailVerified,
		TwoFactorEnabled: u.TwoFactorEnabled,
		Avatar:           u.Avatar,
		Bio:              u.Bio,
		Roles:            u.Roles,
		LastLoginAt:      u.LastLoginAt,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
}

type UserResponse struct {
	ID               uuid.UUID  `json:"id"`
	Email            string     `json:"email"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Phone            string     `json:"phone"`
	Role             string     `json:"role"`
	Status           string     `json:"status"`
	EmailVerified    bool       `json:"email_verified"`
	TwoFactorEnabled bool       `json:"two_factor_enabled"`
	Avatar           string     `json:"avatar"`
	Bio              string     `json:"bio"`
	Roles            []string   `json:"roles"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type CreateUserInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"omitempty"`
	Role      string `json:"role"`
	TenantID  string `json:"tenant_id"`
}

type UpdateUserInput struct {
	Email            *string `json:"email,omitempty"`
	FirstName        *string `json:"first_name,omitempty"`
	LastName         *string `json:"last_name,omitempty"`
	Phone            *string `json:"phone,omitempty"`
	Role             *string `json:"role,omitempty"`
	Status           *string `json:"status,omitempty"`
	Avatar           *string `json:"avatar,omitempty"`
	Bio              *string `json:"bio,omitempty"`
	TwoFactorEnabled *bool   `json:"two_factor_enabled,omitempty"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type UserFilter struct {
	Email    string `query:"email"`
	Role     string `query:"role"`
	Status   string `query:"status"`
	TenantID string `query:"tenant_id"`
	Search   string `query:"search"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}
