package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestErrorCode_Category(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		expected string
	}{
		{ErrCodeValidation, "validation"},
		{ErrCodeInvalidInput, "validation"},
		{ErrCodeNetwork, "network"},
		{ErrCodeConnectionFailed, "network"},
		{ErrCodeDatabase, "database"},
		{ErrCodeQueryFailed, "database"},
		{ErrCodeAuth, "authentication"},
		{ErrCodeUnauthorized, "authentication"},
		{ErrCodeInternal, "internal"},
		{ErrCodePanic, "internal"},
		{ErrCodeFileSystem, "filesystem"},
		{ErrCodeFileNotFound, "filesystem"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.code.Category(); got != tt.expected {
				t.Errorf("ErrorCode.Category() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	err := New(ErrCodeInvalidInput, "invalid input data")

	if err.Code != ErrCodeInvalidInput {
		t.Errorf("New() error code = %v, want %v", err.Code, ErrCodeInvalidInput)
	}
	if err.Message != "invalid input data" {
		t.Errorf("New() message = %v, want %v", err.Message, "invalid input data")
	}
	if err.Category != "validation" {
		t.Errorf("New() category = %v, want %v", err.Category, "validation")
	}
}

func TestError_Error(t *testing.T) {
	inner := errors.New("inner error")
	err := Wrap(inner, ErrCodeQueryFailed, "database operation failed")

	expected := "[3011] database operation failed: inner error"
	if err.Error() != expected {
		t.Errorf("Error.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestError_Unwrap(t *testing.T) {
	inner := errors.New("original error")
	err := Wrap(inner, ErrCodeNetwork, "connection failed")

	if !Is(err.Err, inner) {
		t.Error("Error.Unwrap() should return the wrapped error")
	}
}

func TestError_WithDetails(t *testing.T) {
	err := New(ErrCodeInvalidInput, "validation failed").
		WithDetails("field 'email' is invalid")

	if err.Details != "field 'email' is invalid" {
		t.Errorf("Error.Details = %v, want %v", err.Details, "field 'email' is invalid")
	}
}

func TestError_WithMetadata(t *testing.T) {
	err := New(ErrCodeNetwork, "request failed").
		WithMetadata("url", "https://api.example.com").
		WithMetadata("status", 500)

	if err.Metadata["url"] != "https://api.example.com" {
		t.Errorf("Error.Metadata[url] = %v", err.Metadata["url"])
	}
	if err.Metadata["status"] != 500 {
		t.Errorf("Error.Metadata[status] = %v", err.Metadata["status"])
	}
}

func TestError_WithStack(t *testing.T) {
	err := New(ErrCodeInternal, "internal error").WithStack()

	if err.Stack == "" {
		t.Error("Error.WithStack() should set Stack field")
	}
}

func TestError_ToJSON(t *testing.T) {
	err := New(ErrCodeValidation, "validation error").
		WithDetails("email is required").
		WithMetadata("field", "email")

	data, jsonErr := err.ToJSON()
	if jsonErr != nil {
		t.Fatalf("Error.ToJSON() error = %v", jsonErr)
	}

	var parsed Error
	if parseErr := json.Unmarshal(data, &parsed); parseErr != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", parseErr)
	}

	if parsed.Code != ErrCodeValidation {
		t.Errorf("parsed.Code = %v, want %v", parsed.Code, ErrCodeValidation)
	}
	if parsed.Message != "validation error" {
		t.Errorf("parsed.Message = %v, want %v", parsed.Message, "validation error")
	}
	if parsed.Details != "email is required" {
		t.Errorf("parsed.Details = %v, want %v", parsed.Details, "email is required")
	}
}

func TestFromJSON(t *testing.T) {
	jsonData := []byte(`{
		"code": 1001,
		"message": "invalid input",
		"category": "validation",
		"details": "field required"
	}`)

	err, parseErr := FromJSON(jsonData)
	if parseErr != nil {
		t.Fatalf("FromJSON() error = %v", parseErr)
	}

	if err.Code != ErrCodeInvalidInput {
		t.Errorf("err.Code = %v, want %v", err.Code, ErrCodeInvalidInput)
	}
	if err.Message != "invalid input" {
		t.Errorf("err.Message = %v, want %v", err.Message, "invalid input")
	}
}

func TestWrap(t *testing.T) {
	inner := errors.New("original error")
	err := Wrap(inner, ErrCodeDatabase, "database query failed")

	if err == nil {
		t.Fatal("Wrap() should return non-nil error")
	}
	if err.Code != ErrCodeDatabase {
		t.Errorf("err.Code = %v, want %v", err.Code, ErrCodeDatabase)
	}
	if err.Err != inner {
		t.Error("err.Err should be the original error")
	}
}

func TestWrapWithDetails(t *testing.T) {
	inner := errors.New("connection refused")
	err := WrapWithDetails(inner, ErrCodeConnectionFailed, "connection failed", "timeout after 30s")

	if err.Details != "timeout after 30s" {
		t.Errorf("err.Details = %v, want %v", err.Details, "timeout after 30s")
	}
}

func TestIs(t *testing.T) {
	sentinel := errors.New("sentinel error")
	err := Wrap(sentinel, ErrCodeInternal, "wrapped sentinel")

	if !Is(err, sentinel) {
		t.Error("Is() should return true for sentinel error")
	}
}

func TestAs(t *testing.T) {
	err := New(ErrCodeValidation, "validation error")
	var target *Error
	if !As(err, &target) {
		t.Error("As() should find *Error type")
	}
	if target.Code != ErrCodeValidation {
		t.Errorf("target.Code = %v, want %v", target.Code, ErrCodeValidation)
	}
}

func TestAggregateError_New(t *testing.T) {
	agg := NewAggregateError("multiple validation errors", ErrCodeValidation)

	if agg.Code != ErrCodeValidation {
		t.Errorf("AggregateError.Code = %v, want %v", agg.Code, ErrCodeValidation)
	}
	if agg.Message != "multiple validation errors" {
		t.Errorf("AggregateError.Message = %v, want %v", agg.Message, "multiple validation errors")
	}
	if agg.Count != 0 {
		t.Errorf("AggregateError.Count = %v, want 0", agg.Count)
	}
}

func TestAggregateError_Add(t *testing.T) {
	agg := NewAggregateError("validation errors", ErrCodeValidation)
	agg.Add(errors.New("error 1"))
	agg.Add(errors.New("error 2"))

	if agg.Count != 2 {
		t.Errorf("AggregateError.Count = %v, want 2", agg.Count)
	}
}

func TestAggregateError_AddIf(t *testing.T) {
	agg := NewAggregateError("validation errors", ErrCodeValidation)
	agg.AddIf(true, errors.New("error when true"))
	agg.AddIf(false, errors.New("error when false"))

	if agg.Count != 1 {
		t.Errorf("AggregateError.Count = %v, want 1", agg.Count)
	}
}

func TestAggregateError_HasErrors(t *testing.T) {
	agg := NewAggregateError("errors", ErrCodeValidation)

	if agg.HasErrors() {
		t.Error("HasErrors() should return false for empty aggregate")
	}

	agg.Add(errors.New("error"))
	if !agg.HasErrors() {
		t.Error("HasErrors() should return true when errors exist")
	}
}

func TestAggregateError_Errors(t *testing.T) {
	agg := NewAggregateError("errors", ErrCodeValidation)
	agg.Add(errors.New("error 1"))
	agg.Add(errors.New("error 2"))

	errs := agg.Errors_()
	if len(errs) != 2 {
		t.Errorf("len(Errors_()) = %v, want 2", len(errs))
	}
}

func TestAggregateError_Error(t *testing.T) {
	agg := NewAggregateError("validation failed", ErrCodeValidation)
	agg.Add(errors.New("error 1"))
	agg.Add(errors.New("error 2"))

	errMsg := agg.Error()
	if len(errMsg) == 0 {
		t.Error("AggregateError.Error() should not be empty")
	}
}

func TestAggregateError_ToJSON(t *testing.T) {
	agg := NewAggregateError("validation errors", ErrCodeValidation)
	agg.Add(errors.New("error 1"))

	data, err := agg.ToJSON()
	if err != nil {
		t.Fatalf("AggregateError.ToJSON() error = %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if parsed["code"] != float64(ErrCodeValidation) {
		t.Errorf("parsed code mismatch")
	}
	if parsed["count"] != float64(1) {
		t.Errorf("parsed count mismatch")
	}
}

func TestRecover(t *testing.T) {
	called := false
	recoverFn := func(panicVal interface{}) error {
		called = true
		return New(ErrCodePanic, fmt.Sprintf("panic recovered: %v", panicVal)).WithStack()
	}

	func() {
		defer Recover(recoverFn)
		panic("test panic")
	}()

	if !called {
		t.Error("Recover() should call the recovery function")
	}
}

func TestSafeCall(t *testing.T) {
	fn := func() error {
		return errors.New("function error")
	}

	err := SafeCall(fn)
	if err == nil {
		t.Error("SafeCall() should return error")
	}
}

func TestSafeCall_Panic(t *testing.T) {
	fn := func() error {
		panic("panic in function")
	}

	err := SafeCall(fn, PanicHandler)
	if err == nil {
		t.Error("SafeCall() should return error on panic")
	}
}

func TestSafeCall_Success(t *testing.T) {
	fn := func() error {
		return nil
	}

	err := SafeCall(fn)
	if err != nil {
		t.Errorf("SafeCall() error = %v, want nil", err)
	}
}

func TestCaptureStackTrace(t *testing.T) {
	stack := CaptureStackTrace(1)

	if len(stack) == 0 {
		t.Error("CaptureStackTrace() should return non-empty stack")
	}
}

func TestStackTrace_String(t *testing.T) {
	stack := CaptureStackTrace(1)
	str := stack.String()

	if len(str) == 0 {
		t.Error("StackTrace.String() should return non-empty string")
	}
}

func TestStackTrace_ToJSON(t *testing.T) {
	stack := CaptureStackTrace(1)
	data, err := stack.ToJSON()

	if err != nil {
		t.Fatalf("StackTrace.ToJSON() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("StackTrace.ToJSON() should return non-empty data")
	}
}

func TestErrorFormatter_Format(t *testing.T) {
	err := New(ErrCodeValidation, "validation error").
		WithDetails("email is invalid").
		WithMetadata("field", "email")

	formatter := NewErrorFormatter(WithStack(true), WithCode(true), WithVerbose(true))
	output := formatter.Format(err)

	if len(output) == 0 {
		t.Error("ErrorFormatter.Format() should return non-empty output")
	}
}

func TestErrorFormatter_Default(t *testing.T) {
	err := New(ErrCodeNetwork, "connection failed")
	formatter := NewErrorFormatter()
	output := formatter.Format(err)

	if len(output) == 0 {
		t.Error("ErrorFormatter.Format() should return non-empty output")
	}
}

func TestErrorFormatter_WithoutCode(t *testing.T) {
	err := New(ErrCodeDatabase, "query failed")
	formatter := NewErrorFormatter(WithCode(false))
	output := formatter.Format(err)

	if output == "[3000] query failed" {
		t.Error("Formatter should not include code when disabled")
	}
}

func TestErrorFormatter_WithoutStack(t *testing.T) {
	err := New(ErrCodeInternal, "internal error").WithStack()
	formatter := NewErrorFormatter(WithStack(false))
	output := formatter.Format(err)

	if err.Stack != "" && output == "" {
		t.Error("Formatter should handle missing stack gracefully")
	}
}

func TestError_WithCaller(t *testing.T) {
	err := New(ErrCodeInternal, "internal error").WithCaller(1)

	if err.Stack == "" {
		t.Error("Error.WithCaller() should set Stack field")
	}
}
