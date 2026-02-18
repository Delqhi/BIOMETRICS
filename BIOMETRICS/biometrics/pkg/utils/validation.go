package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Validator struct {
	errors map[string][]string
}

func NewValidator() *Validator {
	return &Validator{
		errors: make(map[string][]string),
	}
}

func (v *Validator) AddError(field, message string) {
	v.errors[field] = append(v.errors[field], message)
}

func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

func (v *Validator) Errors() map[string][]string {
	return v.errors
}

func (v *Validator) Error() string {
	var errs []string
	for field, messages := range v.errors {
		errs = append(errs, fmt.Sprintf("%s: %s", field, strings.Join(messages, ", ")))
	}
	return strings.Join(errs, "; ")
}

func (v *Validator) ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		v.AddError("email", "invalid email format")
		return false
	}
	return true
}

func (v *Validator) ValidateRequired(value, fieldName string) bool {
	if strings.TrimSpace(value) == "" {
		v.AddError(fieldName, "field is required")
		return false
	}
	return true
}

func (v *Validator) ValidateMinLength(value, fieldName string, min int) bool {
	if len(value) < min {
		v.AddError(fieldName, fmt.Sprintf("minimum length is %d", min))
		return false
	}
	return true
}

func (v *Validator) ValidateMaxLength(value, fieldName string, max int) bool {
	if len(value) > max {
		v.AddError(fieldName, fmt.Sprintf("maximum length is %d", max))
		return false
	}
	return true
}

func (v *Validator) ValidateRange(value, fieldName string, min, max int) bool {
	length := len(value)
	if length < min || length > max {
		v.AddError(fieldName, fmt.Sprintf("length must be between %d and %d", min, max))
		return false
	}
	return true
}

func (v *Validator) ValidatePattern(value, fieldName, pattern string) bool {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		v.AddError(fieldName, "invalid pattern")
		return false
	}
	if !matched {
		v.AddError(fieldName, "does not match required pattern")
		return false
	}
	return true
}

func (v *Validator) ValidateURL(url string) bool {
	pattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, _ := regexp.MatchString(pattern, url)
	if !matched {
		v.AddError("url", "invalid URL format")
		return false
	}
	return true
}

func (v *Validator) ValidatePassword(password string) bool {
	if len(password) < 8 {
		v.AddError("password", "password must be at least 8 characters")
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		v.AddError("password", "password must contain uppercase, lowercase, digit, and special character")
		return false
	}
	return true
}

func ValidateStruct(obj interface{}) map[string][]string {
	v := NewValidator()
	validateValue(v, reflect.ValueOf(obj), "")
	return v.errors
}

func validateValue(v *Validator, val reflect.Value, prefix string) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := val.Field(i)

		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fieldName := field.Name
		if prefix != "" {
			fieldName = prefix + "." + fieldName
		}

		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if rule == "required" {
				if fieldValue.Kind() == reflect.String {
					if fieldValue.String() == "" {
						v.AddError(fieldName, "is required")
					}
				}
			}
		}

		if fieldValue.Kind() == reflect.Struct {
			validateValue(v, fieldValue, fieldName)
		}
	}
}

func IsValidUUID(s string) bool {
	pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	return regexp.MustCompile(pattern).MatchString(strings.ToLower(s))
}

func IsValidPhoneNumber(phone string) bool {
	pattern := `^\+?[1-9]\d{1,14}$`
	return regexp.MustCompile(pattern).MatchString(phone)
}
