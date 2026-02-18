package models

import (
	"time"

	"github.com/google/uuid"
)

type Biometric struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	UserID       uuid.UUID              `json:"user_id" db:"user_id"`
	TenantID     *uuid.UUID             `json:"tenant_id,omitempty" db:"tenant_id"`
	Type         string                 `json:"type" db:"type"`
	SubType      string                 `json:"sub_type" db:"sub_type"`
	Label        string                 `json:"label" db:"label"`
	Template     string                 `json:"-" db:"template"`
	Algorithm    string                 `json:"algorithm" db:"algorithm"`
	DeviceID     string                 `json:"device_id" db:"device_id"`
	DeviceInfo   string                 `json:"device_info" db:"device_info"`
	PublicKey    string                 `json:"public_key" db:"public_key"`
	CredentialID string                 `json:"credential_id" db:"credential_id"`
	Status       string                 `json:"status" db:"status"`
	EnrolledAt   time.Time              `json:"enrolled_at" db:"enrolled_at"`
	LastUsedAt   *time.Time             `json:"last_used_at" db:"last_used_at"`
	ExpiresAt    *time.Time             `json:"expires_at" db:"expires_at"`
	FailureCount int                    `json:"failure_count" db:"failure_count"`
	MaxFailures  int                    `json:"max_failures" db:"max_failures"`
	VerifyCount  int                    `json:"verify_count" db:"verify_count"`
	SuccessCount int                    `json:"success_count" db:"success_count"`
	QualityScore int                    `json:"quality_score" db:"quality_score"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	IsActive     bool                   `json:"is_active" db:"is_active"`
	IsPrimary    bool                   `json:"is_primary" db:"is_primary"`
	Version      int                    `json:"version" db:"version"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
}

type BiometricType string

const (
	BiometricTypeFingerprint BiometricType = "fingerprint"
	BiometricTypeFace        BiometricType = "face"
	BiometricTypeIris        BiometricType = "iris"
	BiometricTypeVoice       BiometricType = "voice"
	BiometricTypePalm        BiometricType = "palm"
	BiometricTypeRetina      BiometricType = "retina"
)

type BiometricSubType string

const (
	BiometricSubTypeLeftThumb  BiometricSubType = "left_thumb"
	BiometricSubTypeRightThumb BiometricSubType = "right_thumb"
	BiometricSubTypeLeftIndex  BiometricSubType = "left_index"
	BiometricSubTypeRightIndex BiometricSubType = "right_index"
	BiometricSubType3D         BiometricSubType = "3d"
	BiometricSubType2D         BiometricSubType = "2d"
)

type BiometricStatus string

const (
	BiometricStatusEnrolled  BiometricStatus = "enrolled"
	BiometricStatusActive    BiometricStatus = "active"
	BiometricStatusSuspended BiometricStatus = "suspended"
	BiometricStatusRevoked   BiometricStatus = "revoked"
	BiometricStatusExpired   BiometricStatus = "expired"
)

type BiometricAlgorithm string

const (
	AlgorithmFIDO2       BiometricAlgorithm = "fido2"
	AlgorithmBioKey      BiometricAlgorithm = "biokey"
	AlgorithmProprietary BiometricAlgorithm = "proprietary"
)

type BiometricVerification struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	BiometricID   uuid.UUID  `json:"biometric_id" db:"biometric_id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	TenantID      *uuid.UUID `json:"tenant_id,omitempty" db:"tenant_id"`
	Challenge     string     `json:"challenge" db:"challenge"`
	ChallengeHash string     `json:"-" db:"challenge_hash"`
	ClientData    string     `json:"client_data" db:"client_data"`
	AuthData      string     `json:"auth_data" db:"auth_data"`
	Signature     string     `json:"-" db:"signature"`
	Result        string     `json:"result" db:"result"`
	Score         int        `json:"score" db:"score"`
	FailureReason string     `json:"failure_reason" db:"failure_reason"`
	IPAddress     string     `json:"ip_address" db:"ip_address"`
	UserAgent     string     `json:"user_agent" db:"user_agent"`
	DeviceInfo    string     `json:"device_info" db:"device_info"`
	Location      string     `json:"location" db:"location"`
	Latency       int        `json:"latency" db:"latency"`
	VerifiedAt    time.Time  `json:"verified_at" db:"verified_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

type VerificationResult string

const (
	VerificationResultSuccess   VerificationResult = "success"
	VerificationResultFailed    VerificationResult = "failed"
	VerificationResultPending   VerificationResult = "pending"
	VerificationResultCancelled VerificationResult = "cancelled"
	VerificationResultTimeout   VerificationResult = "timeout"
)

type CreateBiometricInput struct {
	UserID       string `json:"user_id" binding:"required"`
	TenantID     string `json:"tenant_id"`
	Type         string `json:"type" binding:"required"`
	SubType      string `json:"sub_type"`
	Label        string `json:"label"`
	DeviceID     string `json:"device_id"`
	PublicKey    string `json:"public_key"`
	CredentialID string `json:"credential_id"`
	QualityScore int    `json:"quality_score"`
}

type BiometricFilter struct {
	Type     string `query:"type"`
	SubType  string `query:"sub_type"`
	UserID   string `query:"user_id"`
	Status   string `query:"status"`
	TenantID string `query:"tenant_id"`
	DeviceID string `query:"device_id"`
	IsActive *bool  `query:"is_active"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}

type BiometricResponse struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	TenantID     *uuid.UUID `json:"tenant_id,omitempty"`
	Type         string     `json:"type"`
	SubType      string     `json:"sub_type"`
	Label        string     `json:"label"`
	DeviceID     string     `json:"device_id"`
	DeviceInfo   string     `json:"device_info"`
	Status       string     `json:"status"`
	EnrolledAt   time.Time  `json:"enrolled_at"`
	LastUsedAt   *time.Time `json:"last_used_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	FailureCount int        `json:"failure_count"`
	VerifyCount  int        `json:"verify_count"`
	SuccessCount int        `json:"success_count"`
	QualityScore int        `json:"quality_score"`
	IsActive     bool       `json:"is_active"`
	IsPrimary    bool       `json:"is_primary"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (b *Biometric) ToResponse() *BiometricResponse {
	return &BiometricResponse{
		ID:           b.ID,
		UserID:       b.UserID,
		TenantID:     b.TenantID,
		Type:         b.Type,
		SubType:      b.SubType,
		Label:        b.Label,
		DeviceID:     b.DeviceID,
		DeviceInfo:   b.DeviceInfo,
		Status:       b.Status,
		EnrolledAt:   b.EnrolledAt,
		LastUsedAt:   b.LastUsedAt,
		ExpiresAt:    b.ExpiresAt,
		FailureCount: b.FailureCount,
		VerifyCount:  b.VerifyCount,
		SuccessCount: b.SuccessCount,
		QualityScore: b.QualityScore,
		IsActive:     b.IsActive,
		IsPrimary:    b.IsPrimary,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}
