package utils

import (
	"errors"
	"fmt"
	"strings"
)

type AppError struct {
	Code       string
	Message    string
	Err        error
	StackTrace string
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func WrapError(err error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func IsErrorType(err error, code string) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

type ErrorCode struct {
	Code        string
	Message     string
	HTTPStatus  int
	IsRetryable bool
}

var (
	ErrNotFound           = ErrorCode{"NOT_FOUND", "Resource not found", 404, false}
	ErrUnauthorized       = ErrorCode{"UNAUTHORIZED", "Unauthorized access", 401, false}
	ErrForbidden          = ErrorCode{"FORBIDDEN", "Forbidden access", 403, false}
	ErrBadRequest         = ErrorCode{"BAD_REQUEST", "Bad request", 400, false}
	ErrInternalServer     = ErrorCode{"INTERNAL_SERVER", "Internal server error", 500, false}
	ErrConflict           = ErrorCode{"CONFLICT", "Resource conflict", 409, false}
	ErrTooManyRequests    = ErrorCode{"TOO_MANY_REQUESTS", "Too many requests", 429, true}
	ErrServiceUnavailable = ErrorCode{"SERVICE_UNAVAILABLE", "Service unavailable", 503, true}
	ErrTimeout            = ErrorCode{"TIMEOUT", "Request timeout", 504, true}
)

func (ec ErrorCode) Error() string {
	return ec.Message
}

func (ec ErrorCode) WithMessage(msg string) *AppError {
	return &AppError{
		Code:    ec.Code,
		Message: msg,
	}
}

func (ec ErrorCode) WithErr(err error) *AppError {
	return &AppError{
		Code:    ec.Code,
		Message: ec.Message,
		Err:     err,
	}
}

func ToMap(err error) map[string]interface{} {
	result := map[string]interface{}{
		"error": err.Error(),
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		result["code"] = appErr.Code
		result["message"] = appErr.Message
		if appErr.Err != nil {
			result["details"] = appErr.Err.Error()
		}
	}

	return result
}

func ChainErrors(errs ...error) error {
	var errStrings []string
	for _, err := range errs {
		if err != nil {
			errStrings = append(errStrings, err.Error())
		}
	}
	if len(errStrings) == 0 {
		return nil
	}
	if len(errStrings) == 1 {
		return errors.New(errStrings[0])
	}
	return errors.New(strings.Join(errStrings, "; "))
}
