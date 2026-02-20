// Package errors provides comprehensive error handling utilities for the biometrics CLI.
//
// This package implements:
//   - Error codes with categories (validation, network, database, auth, etc.)
//   - Structured error types with stack traces
//   - Error wrapping and unwrapping
//   - Recovery utilities with panic handling
//   - Error aggregation for multiple errors
//   - JSON serialization for API responses
//
// Best Practices Feb 2026 compliant.
package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
)

// ErrorCode represents a unique error code for categorization
type ErrorCode int

// Error categories
const (
	// Validation errors (1000-1999)
	ErrCodeValidation ErrorCode = 1000 + iota
	ErrCodeInvalidInput
	ErrCodeMissingField
	ErrCodeMalformedData
	ErrCodeConstraintViolation

	// Network errors (2000-2999)
	ErrCodeNetwork ErrorCode = 2000 + iota
	ErrCodeConnectionFailed
	ErrCodeTimeout
	ErrCodeTooManyRequests
	ErrCodeServiceUnavailable

	// Database errors (3000-3999)
	ErrCodeDatabase ErrorCode = 3000 + iota
	ErrCodeQueryFailed
	ErrCodeConnectionLost
	ErrCodeTransactionFailed
	ErrCodeDataIntegrity

	// Authentication/Authorization errors (4000-4999)
	ErrCodeAuth ErrorCode = 4000 + iota
	ErrCodeUnauthorized
	ErrCodeForbidden
	ErrCodeTokenExpired
	ErrCodeInvalidCredentials

	// Internal errors (5000-5999)
	ErrCodeInternal ErrorCode = 5000 + iota
	ErrCodePanic
	ErrCodeNotImplemented
	ErrCodeConfiguration
	ErrCodeDependency

	// File system errors (6000-6999)
	ErrCodeFileSystem ErrorCode = 6000 + iota
	ErrCodeFileNotFound
	ErrCodePermissionDenied
	ErrCodeDiskFull
)

// ErrorCategory returns the category string for an error code
func (e ErrorCode) Category() string {
	switch {
	case e >= 1000 && e < 2000:
		return "validation"
	case e >= 2000 && e < 3000:
		return "network"
	case e >= 3000 && e < 4000:
		return "database"
	case e >= 4000 && e < 5000:
		return "authentication"
	case e >= 5000 && e < 6000:
		return "internal"
	case e >= 6000 && e < 7000:
		return "filesystem"
	default:
		return "unknown"
	}
}

// Error represents a structured error with additional context
type Error struct {
	Code     ErrorCode              `json:"code"`
	Message  string                 `json:"message"`
	Category string                 `json:"category"`
	Details  string                 `json:"details,omitempty"`
	Stack    string                 `json:"stack,omitempty"`
	Err      error                  `json:"-"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	mu       sync.RWMutex
}

// Error returns the error message
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// WithDetails adds additional details to the error
func (e *Error) WithDetails(details string) *Error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Details = details
	return e
}

// WithMetadata adds metadata to the error
func (e *Error) WithMetadata(key string, value interface{}) *Error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WithStack adds stack trace to the error
func (e *Error) WithStack() *Error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, file, line, _ := runtime.Caller(1)
	e.Stack = fmt.Sprintf("%s:%d", file, line)
	return e
}

// WithCaller adds specific caller information to the error
func (e *Error) WithCaller(depth int) *Error {
	e.mu.Lock()
	defer e.mu.Unlock()
	pc, file, line, _ := runtime.Caller(depth)
	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name()
	}
	e.Stack = fmt.Sprintf("%s:%d (%s)", file, line, funcName)
	return e
}

// ToJSON converts the error to JSON
func (e *Error) ToJSON() ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return json.Marshal(e)
}

// FromJSON creates an error from JSON
func FromJSON(data []byte) (*Error, error) {
	var err Error
	if err := json.Unmarshal(data, &err); err != nil {
		return nil, fmt.Errorf("failed to unmarshal error: %w", err)
	}
	return &err, nil
}

// New creates a new structured error
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		Category: code.Category(),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Code:     code,
		Message:  message,
		Category: code.Category(),
		Err:      err,
	}
}

// WrapWithDetails wraps an error with details
func WrapWithDetails(err error, code ErrorCode, message, details string) *Error {
	e := Wrap(err, code, message)
	e.Details = details
	return e
}

// Is checks if the error matches the target
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in the chain that matches the target
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// AggregateError holds multiple errors
type AggregateError struct {
	Errors  []error   `json:"errors"`
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Count   int       `json:"count"`
	mu      sync.RWMutex
}

// NewAggregateError creates a new aggregate error
func NewAggregateError(message string, code ErrorCode) *AggregateError {
	return &AggregateError{
		Code:    code,
		Message: message,
		Errors:  make([]error, 0),
	}
}

// Add adds an error to the aggregate
func (a *AggregateError) Add(err error) {
	if err == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Errors = append(a.Errors, err)
	a.Count = len(a.Errors)
}

// AddIf adds an error only if condition is true
func (a *AggregateError) AddIf(condition bool, err error) {
	if condition {
		a.Add(err)
	}
}

// Error returns the aggregate error message
func (a *AggregateError) Error() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.Errors) == 0 {
		return "no errors"
	}
	if len(a.Errors) == 1 {
		return a.Errors[0].Error()
	}
	var msgs []string
	for _, e := range a.Errors {
		msgs = append(msgs, e.Error())
	}
	return fmt.Sprintf("%s (%d errors: %s)", a.Message, a.Count, strings.Join(msgs, "; "))
}

// HasErrors returns true if there are any errors
func (a *AggregateError) HasErrors() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.Errors) > 0
}

// Errors returns a copy of the errors slice
func (a *AggregateError) Errors_() []error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]error, len(a.Errors))
	copy(result, a.Errors)
	return result
}

// ErrorCount returns the number of errors
func (a *AggregateError) ErrorCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Count
}

// ToJSON converts aggregate error to JSON
func (a *AggregateError) ToJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a)
}

// RecoverFunc is a function that handles recovered panic values
type RecoverFunc func(panicVal interface{}) error

// Recover recovers from a panic and optionally logs or handles it
func Recover(recoverFn ...RecoverFunc) (err error) {
	if r := recover(); r != nil {
		for _, fn := range recoverFn {
			if e := fn(r); e != nil {
				err = e
			}
		}
		if err == nil {
			err = New(ErrCodePanic, fmt.Sprintf("panic recovered: %v", r)).WithStack()
		}
	}
	return
}

// SafeCall executes a function and recovers from any panics
func SafeCall(fn func() error, recoverFn ...RecoverFunc) (err error) {
	defer func() {
		if r := recover(); r != nil {
			for _, fn := range recoverFn {
				if e := fn(r); e != nil {
					err = e
				}
			}
			if err == nil {
				err = New(ErrCodePanic, fmt.Sprintf("panic in SafeCall: %v", r)).WithStack()
			}
		}
	}()
	return fn()
}

// PanicHandler is a function that handles panic values
func PanicHandler(panicVal interface{}) error {
	return New(ErrCodePanic, fmt.Sprintf("panic occurred: %v", panicVal)).WithStack()
}

// StackTrace represents a stack trace
type StackTrace []StackFrame

// StackFrame represents a single frame in a stack trace
type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// CaptureStackTrace captures the current stack trace
func CaptureStackTrace(skip int) StackTrace {
	stack := make(StackTrace, 0)
	for i := skip; i < skip+10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		funcName := "unknown"
		if fn != nil {
			funcName = fn.Name()
		}
		stack = append(stack, StackFrame{
			Function: funcName,
			File:     file,
			Line:     line,
		})
	}
	return stack
}

// String returns a string representation of the stack trace
func (s StackTrace) String() string {
	var b strings.Builder
	for _, frame := range s {
		fmt.Fprintf(&b, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
	}
	return b.String()
}

// ToJSON converts stack trace to JSON
func (s StackTrace) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// ErrorFormatter provides consistent error formatting
type ErrorFormatter struct {
	includeStack bool
	includeCode  bool
	verbose      bool
}

// NewErrorFormatter creates a new ErrorFormatter
func NewErrorFormatter(opts ...Option) *ErrorFormatter {
	f := &ErrorFormatter{
		includeStack: true,
		includeCode:  true,
		verbose:      false,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// Option is a functional option for ErrorFormatter
type Option func(*ErrorFormatter)

// WithStack includes stack trace in formatted output
func WithStack(include bool) Option {
	return func(f *ErrorFormatter) {
		f.includeStack = include
	}
}

// WithCode includes error code in formatted output
func WithCode(include bool) Option {
	return func(f *ErrorFormatter) {
		f.includeCode = include
	}
}

// WithVerbose enables verbose output
func WithVerbose(verbose bool) Option {
	return func(f *ErrorFormatter) {
		f.verbose = verbose
	}
}

// Format formats an error according to the formatter options
func (f *ErrorFormatter) Format(err error) string {
	var b strings.Builder

	if biErr, ok := err.(*Error); ok {
		if f.includeCode {
			b.WriteString(fmt.Sprintf("[%d] ", biErr.Code))
		}
		b.WriteString(biErr.Message)
		if biErr.Details != "" {
			b.WriteString(fmt.Sprintf(": %s", biErr.Details))
		}
		if f.includeStack && biErr.Stack != "" {
			b.WriteString(fmt.Sprintf("\nStack: %s", biErr.Stack))
		}
		if f.verbose && biErr.Metadata != nil {
			b.WriteString("\nMetadata:")
			for k, v := range biErr.Metadata {
				fmt.Fprintf(&b, "\n  %s: %v", k, v)
			}
		}
	} else {
		b.WriteString(err.Error())
	}

	return b.String()
}
