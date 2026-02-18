package models

import (
	"errors"
	"net/http"

	"golang.org/x/text/message"
)

var (
	ErrNotFound             = errors.New("resource not found")
	ErrAlreadyExists        = errors.New("resource already exists")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidInput         = errors.New("invalid input")
	ErrValidationFailed     = errors.New("validation failed")
	ErrInternalServerError  = errors.New("internal server error")
	ErrServiceUnavailable   = errors.New("service unavailable")
	ErrConflict             = errors.New("conflict")
	ErrGone                 = errors.New("resource gone")
	ErrRateLimited          = errors.New("rate limited")
	ErrTooManyRequests      = errors.New("too many requests")
	ErrExpired              = errors.New("token expired")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenRevoked         = errors.New("token revoked")
	ErrInsufficientScope    = errors.New("insufficient scope")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrAccountLocked        = errors.New("account locked")
	ErrAccountSuspended     = errors.New("account suspended")
	ErrEmailNotVerified     = errors.New("email not verified")
	ErrTwoFactorRequired    = errors.New("two factor authentication required")
	ErrTwoFactorInvalid     = errors.New("invalid two factor code")
	ErrPasswordTooWeak      = errors.New("password too weak")
	ErrPasswordCompromised  = errors.New("password compromised")
	ErrSessionExpired       = errors.New("session expired")
	ErrSessionRevoked       = errors.New("session revoked")
	ErrBiometricNotEnrolled = errors.New("biometric not enrolled")
	ErrBiometricInvalid     = errors.New("invalid biometric")
	ErrBiometricExpired     = errors.New("biometric expired")
	ErrQuotaExceeded        = errors.New("quota exceeded")
	ErrWorkflowNotFound     = errors.New("workflow not found")
	ErrWorkflowFailed       = errors.New("workflow failed")
	ErrIntegrationNotFound  = errors.New("integration not found")
	ErrIntegrationFailed    = errors.New("integration failed")
	ErrContentNotFound      = errors.New("content not found")
	ErrFileNotFound         = errors.New("file not found")
	ErrFileTooLarge         = errors.New("file too large")
	ErrInvalidFileType      = errors.New("invalid file type")
)

type AppError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Details    string            `json:"details,omitempty"`
	StatusCode int               `json:"status_code"`
	Err        error             `json:"-"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return e.Code + ": " + e.Message + " - " + e.Details
	}
	return e.Code + ": " + e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

func NewNotFoundError(resource string) *AppError {
	return NewAppError(
		"NOT_FOUND",
		resource+" not found",
		http.StatusNotFound,
		ErrNotFound,
	)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(
		"UNAUTHORIZED",
		message,
		http.StatusUnauthorized,
		ErrUnauthorized,
	)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(
		"FORBIDDEN",
		message,
		http.StatusForbidden,
		ErrForbidden,
	)
}

func NewValidationError(details string) *AppError {
	return NewAppError(
		"VALIDATION_ERROR",
		"Validation failed",
		http.StatusBadRequest,
		ErrValidationFailed,
	).WithDetails(details)
}

func NewConflictError(resource string) *AppError {
	return NewAppError(
		"CONFLICT",
		resource+" already exists",
		http.StatusConflict,
		ErrAlreadyExists,
	)
}

func NewInternalError(message string, err error) *AppError {
	return NewAppError(
		"INTERNAL_ERROR",
		message,
		http.StatusInternalServerError,
		err,
	)
}

func NewRateLimitError(retryAfter int) *AppError {
	return &AppError{
		Code:       "RATE_LIMITED",
		Message:    "Too many requests",
		StatusCode: http.StatusTooManyRequests,
		Err:        ErrRateLimited,
		Metadata:   map[string]string{"retry_after": string(rune(retryAfter))},
	}
}

func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

func (e *AppError) WithMetadata(key, value string) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

type ErrorResponse struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Details   string            `json:"details,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp string            `json:"timestamp"`
	RequestID string            `json:"request_id,omitempty"`
}

func (e *AppError) ToResponse(requestID string) *ErrorResponse {
	return &ErrorResponse{
		Code:      e.Code,
		Message:   e.Message,
		Details:   e.Details,
		Metadata:  e.Metadata,
		Timestamp: "now",
		RequestID: requestID,
	}
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	return "validation errors"
}

func NewValidationErrors(errors []ValidationError) ValidationErrors {
	return errors
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	HasMore    bool        `json:"has_more"`
}

func NewPaginatedResponse(data interface{}, page, pageSize int, totalItems int64) *PaginatedResponse {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}
	return &PaginatedResponse{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}
}
