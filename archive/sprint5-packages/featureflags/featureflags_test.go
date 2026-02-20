package featureflags

import (
	"os"
	"testing"
)

func TestFeatureFlagManager_New(t *testing.T) {
	m := NewFeatureFlagManager()

	if m == nil {
		t.Error("NewFeatureFlagManager() should not return nil")
	}
	if m.flags == nil {
		t.Error("flags map should be initialized")
	}
	if m.listeners == nil {
		t.Error("listeners map should be initialized")
	}
}

func TestFeatureFlagManager_RegisterBooleanFlag(t *testing.T) {
	m := NewFeatureFlagManager()

	err := m.RegisterBooleanFlag("feature.enabled", true, "Enable feature")
	if err != nil {
		t.Errorf("RegisterBooleanFlag() error = %v", err)
	}

	flag, err := m.GetFlag("feature.enabled")
	if err != nil {
		t.Errorf("GetFlag() error = %v", err)
	}
	if flag.Type != FlagTypeBoolean {
		t.Errorf("flag.Type = %v, want %v", flag.Type, FlagTypeBoolean)
	}
	if flag.Default != true {
		t.Errorf("flag.Default = %v, want true", flag.Default)
	}
}

func TestFeatureFlagManager_RegisterStringFlag(t *testing.T) {
	m := NewFeatureFlagManager()

	err := m.RegisterStringFlag("api.url", "https://api.example.com", "API URL")
	if err != nil {
		t.Errorf("RegisterStringFlag() error = %v", err)
	}

	flag, err := m.GetFlag("api.url")
	if err != nil {
		t.Errorf("GetFlag() error = %v", err)
	}
	if flag.Type != FlagTypeString {
		t.Errorf("flag.Type = %v, want %v", flag.Type, FlagTypeString)
	}
	if flag.Default != "https://api.example.com" {
		t.Errorf("flag.Default = %v, want https://api.example.com", flag.Default)
	}
}

func TestFeatureFlagManager_RegisterIntFlag(t *testing.T) {
	m := NewFeatureFlagManager()

	err := m.RegisterIntFlag("rate.limit", 100, "Rate limit")
	if err != nil {
		t.Errorf("RegisterIntFlag() error = %v", err)
	}

	flag, err := m.GetFlag("rate.limit")
	if err != nil {
		t.Errorf("GetFlag() error = %v", err)
	}
	if flag.Type != FlagTypeInt {
		t.Errorf("flag.Type = %v, want %v", flag.Type, FlagTypeInt)
	}
	if flag.Default != 100 {
		t.Errorf("flag.Default = %v, want 100", flag.Default)
	}
}

func TestFeatureFlagManager_RegisterFloatFlag(t *testing.T) {
	m := NewFeatureFlagManager()

	err := m.RegisterFloatFlag("threshold.value", 0.75, "Threshold value")
	if err != nil {
		t.Errorf("RegisterFloatFlag() error = %v", err)
	}

	flag, err := m.GetFlag("threshold.value")
	if err != nil {
		t.Errorf("GetFlag() error = %v", err)
	}
	if flag.Type != FlagTypeFloat {
		t.Errorf("flag.Type = %v, want %v", flag.Type, FlagTypeFloat)
	}
	if flag.Default != 0.75 {
		t.Errorf("flag.Default = %v, want 0.75", flag.Default)
	}
}

func TestFeatureFlagManager_RegisterJSONFlag(t *testing.T) {
	m := NewFeatureFlagManager()

	err := m.RegisterJSONFlag("config.settings", map[string]interface{}{"key": "value"}, "Config settings")
	if err != nil {
		t.Errorf("RegisterJSONFlag() error = %v", err)
	}

	flag, err := m.GetFlag("config.settings")
	if err != nil {
		t.Errorf("GetFlag() error = %v", err)
	}
	if flag.Type != FlagTypeJSON {
		t.Errorf("flag.Type = %v, want %v", flag.Type, FlagTypeJSON)
	}
}

func TestFeatureFlagManager_GetBool(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterBooleanFlag("feature.enabled", false, "Enable feature")

	if m.GetBool("feature.enabled") != false {
		t.Error("GetBool() should return default value")
	}

	m.SetBool("feature.enabled", true)
	if m.GetBool("feature.enabled") != true {
		t.Error("GetBool() should return set value")
	}
}

func TestFeatureFlagManager_GetString(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterStringFlag("api.url", "default.example.com", "API URL")

	if m.GetString("api.url") != "default.example.com" {
		t.Error("GetString() should return default value")
	}

	m.SetString("api.url", "custom.example.com")
	if m.GetString("api.url") != "custom.example.com" {
		t.Error("GetString() should return set value")
	}
}

func TestFeatureFlagManager_GetInt(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterIntFlag("rate.limit", 50, "Rate limit")

	if m.GetInt("rate.limit") != 50 {
		t.Error("GetInt() should return default value")
	}

	m.SetInt("rate.limit", 100)
	if m.GetInt("rate.limit") != 100 {
		t.Error("GetInt() should return set value")
	}
}

func TestFeatureFlagManager_GetFloat(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterFloatFlag("threshold.value", 0.5, "Threshold")

	if m.GetFloat("threshold.value") != 0.5 {
		t.Error("GetFloat() should return default value")
	}

	m.SetFloat("threshold.value", 0.75)
	if m.GetFloat("threshold.value") != 0.75 {
		t.Error("GetFloat() should return set value")
	}
}

func TestFeatureFlagManager_IsEnabled(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterBooleanFlag("feature.enabled", true, "Enable feature")

	if !m.IsEnabled("feature.enabled") {
		t.Error("IsEnabled() should return true for registered flag")
	}

	m.SetEnabled("feature.enabled", false)
	if m.IsEnabled("feature.enabled") {
		t.Error("IsEnabled() should return false after disabling")
	}
}

func TestFeatureFlagManager_SetEnabled(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterBooleanFlag("feature.enabled", true, "Enable feature")

	err := m.SetEnabled("feature.enabled", false)
	if err != nil {
		t.Errorf("SetEnabled() error = %v", err)
	}

	if m.IsEnabled("feature.enabled") {
		t.Error("IsEnabled() should return false after SetEnabled(false)")
	}
}

func TestFeatureFlagManager_GetFlag_NotFound(t *testing.T) {
	m := NewFeatureFlagManager()

	_, err := m.GetFlag("nonexistent")
	if err == nil {
		t.Error("GetFlag() should return error for nonexistent flag")
	}
}

func TestFeatureFlagManager_GetAllFlags(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterBooleanFlag("flag1", true, "Flag 1")
	m.RegisterStringFlag("flag2", "value", "Flag 2")

	flags := m.GetAllFlags()
	if len(flags) != 2 {
		t.Errorf("len(GetAllFlags()) = %v, want 2", len(flags))
	}
}

func TestFeatureFlagManager_RegisterChangeListener(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterBooleanFlag("feature.enabled", false, "Enable feature")

	var called bool
	m.RegisterChangeListener("feature.enabled", func(name string, oldValue, newValue interface{}) {
		called = true
		if oldValue != false {
			t.Errorf("oldValue = %v, want false", oldValue)
		}
		if newValue != true {
			t.Errorf("newValue = %v, want true", newValue)
		}
	})

	m.SetBool("feature.enabled", true)
	if !called {
		t.Error("ChangeListener should be called on flag change")
	}
}

func TestFeatureFlagManager_RegisterDuplicate(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterBooleanFlag("feature.enabled", true, "Enable feature")

	err := m.RegisterBooleanFlag("feature.enabled", false, "Enable feature again")
	if err == nil {
		t.Error("RegisterDuplicate should return error")
	}
}

func TestFeatureFlagManager_LoadFromEnvironment(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterBooleanFlag("feature.enabled", false, "Enable feature")
	m.RegisterStringFlag("api.url", "default.example.com", "API URL")
	m.RegisterIntFlag("rate.limit", 50, "Rate limit")

	os.Setenv("MYAPP_FEATURE_ENABLED", "true")
	os.Setenv("MYAPP_API_URL", "env.example.com")
	os.Setenv("MYAPP_RATE_LIMIT", "200")
	defer os.Unsetenv("MYAPP_FEATURE_ENABLED")
	defer os.Unsetenv("MYAPP_API_URL")
	defer os.Unsetenv("MYAPP_RATE_LIMIT")

	err := m.LoadFromEnvironment("myapp")
	if err != nil {
		t.Errorf("LoadFromEnvironment() error = %v", err)
	}

	if m.GetBool("feature.enabled") != true {
		t.Error("LoadFromEnvironment should set boolean from env")
	}
	if m.GetString("api.url") != "env.example.com" {
		t.Error("LoadFromEnvironment should set string from env")
	}
	if m.GetInt("rate.limit") != 200 {
		t.Error("LoadFromEnvironment should set int from env")
	}
}

func TestFeatureFlagManager_EvaluateWithContext(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterBooleanFlag("feature.enabled", false, "Enable feature")
	m.RegisterStringFlag("greeting", "Hello {name}", "Greeting template")

	m.SetBool("feature.enabled", true)
	m.SetString("greeting", "Hello {name}")

	result, err := m.EvaluateWithContext("feature.enabled", map[string]interface{}{"name": "World"})
	if err != nil {
		t.Errorf("EvaluateWithContext() error = %v", err)
	}
	if result != true {
		t.Errorf("EvaluateWithContext() = %v, want true", result)
	}

	resultStr, err := m.EvaluateWithContext("greeting", map[string]interface{}{"name": "World"})
	if err != nil {
		t.Errorf("EvaluateWithContext() error = %v", err)
	}
	if resultStr != "Hello World" {
		t.Errorf("EvaluateWithContext() = %v, want 'Hello World'", resultStr)
	}
}

func TestFlagType_Category(t *testing.T) {
	tests := []struct {
		flagType FlagType
		expected string
	}{
		{FlagTypeBoolean, "boolean"},
		{FlagTypeString, "string"},
		{FlagTypeInt, "int"},
		{FlagTypeFloat, "float"},
		{FlagTypeJSON, "json"},
	}

	for _, tt := range tests {
		t.Run(string(tt.flagType), func(t *testing.T) {
			if string(tt.flagType) != tt.expected {
				t.Errorf("FlagType = %v, want %v", tt.flagType, tt.expected)
			}
		})
	}
}

func TestValidateFlagType(t *testing.T) {
	tests := []struct {
		value        interface{}
		expectedType FlagType
		valid        bool
	}{
		{true, FlagTypeBoolean, true},
		{false, FlagTypeBoolean, true},
		{"string", FlagTypeString, true},
		{123, FlagTypeInt, true},
		{int64(123), FlagTypeInt, true},
		{1.5, FlagTypeFloat, true},
		{map[string]interface{}{}, FlagTypeJSON, true},
		{"string", FlagTypeBoolean, false},
		{123, FlagTypeString, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := ValidateFlagType(tt.value, tt.expectedType)
			if result != tt.valid {
				t.Errorf("ValidateFlagType(%v, %v) = %v, want %v", tt.value, tt.expectedType, result, tt.valid)
			}
		})
	}
}

func TestGetFlagType(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected FlagType
	}{
		{true, FlagTypeBoolean},
		{"string", FlagTypeString},
		{123, FlagTypeInt},
		{1.5, FlagTypeFloat},
		{map[string]interface{}{}, FlagTypeJSON},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := GetFlagType(tt.value)
			if result != tt.expected {
				t.Errorf("GetFlagType(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestGetFlagType_Nil(t *testing.T) {
	result := GetFlagType(nil)
	if result != FlagTypeJSON {
		t.Errorf("GetFlagType(nil) = %v, want %v", result, FlagTypeJSON)
	}
}

func TestUnmarshalFlagValue(t *testing.T) {
	m := NewFeatureFlagManager()
	m.RegisterStringFlag("test.flag", "default", "Test flag")
	m.SetString("test.flag", "custom")

	flag, _ := m.GetFlag("test.flag")

	var target string
	err := UnmarshalFlagValue(flag, &target)
	if err != nil {
		t.Errorf("UnmarshalFlagValue() error = %v", err)
	}
	if target != "custom" {
		t.Errorf("target = %v, want 'custom'", target)
	}
}

func TestUnmarshalFlagValue_NilFlag(t *testing.T) {
	var target string
	err := UnmarshalFlagValue(nil, &target)
	if err == nil {
		t.Error("UnmarshalFlagValue(nil) should return error")
	}
}

func TestFeatureFlagManager_GetBool_NotFound(t *testing.T) {
	m := NewFeatureFlagManager()

	result := m.GetBool("nonexistent")
	if result != false {
		t.Errorf("GetBool(nonexistent) = %v, want false", result)
	}
}

func TestFeatureFlagManager_GetString_NotFound(t *testing.T) {
	m := NewFeatureFlagManager()

	result := m.GetString("nonexistent")
	if result != "" {
		t.Errorf("GetString(nonexistent) = %v, want empty string", result)
	}
}

func TestFeatureFlagManager_GetInt_NotFound(t *testing.T) {
	m := NewFeatureFlagManager()

	result := m.GetInt("nonexistent")
	if result != 0 {
		t.Errorf("GetInt(nonexistent) = %v, want 0", result)
	}
}

func TestFeatureFlagManager_GetFloat_NotFound(t *testing.T) {
	m := NewFeatureFlagManager()

	result := m.GetFloat("nonexistent")
	if result != 0.0 {
		t.Errorf("GetFloat(nonexistent) = %v, want 0.0", result)
	}
}

func TestFeatureFlagManager_IsEnabled_NotFound(t *testing.T) {
	m := NewFeatureFlagManager()

	result := m.IsEnabled("nonexistent")
	if result != false {
		t.Errorf("IsEnabled(nonexistent) = %v, want false", result)
	}
}
