// Package featureflags provides runtime feature flag management for the biometrics CLI.
//
// This package implements:
//   - Runtime feature flag management
//   - Boolean, string, int, json flag types
//   - Flag evaluation with context
//   - Default values
//   - Flag change listeners
//   - File-based and environment-based configuration
//
// Best Practices Feb 2026 compliant.
package featureflags

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FlagType represents the type of a feature flag
type FlagType string

const (
	// FlagTypeBoolean represents a boolean flag
	FlagTypeBoolean FlagType = "boolean"
	// FlagTypeString represents a string flag
	FlagTypeString FlagType = "string"
	// FlagTypeInt represents an integer flag
	FlagTypeInt FlagType = "int"
	// FlagTypeFloat represents a float flag
	FlagTypeFloat FlagType = "float"
	// FlagTypeJSON represents a JSON flag
	FlagTypeJSON FlagType = "json"
)

// Flag represents a feature flag with its configuration
type Flag struct {
	Name        string      `json:"name"`
	Type        FlagType    `json:"type"`
	Value       interface{} `json:"value"`
	Default     interface{} `json:"default"`
	Description string      `json:"description,omitempty"`
	Enabled     bool        `json:"enabled"`
	UpdatedAt   time.Time   `json:"updated_at"`
	mu          sync.RWMutex
}

// FlagConfig contains configuration for creating a new flag
type FlagConfig struct {
	Name        string
	Type        FlagType
	Default     interface{}
	Description string
	Enabled     bool
}

// FeatureFlagManager manages all feature flags
type FeatureFlagManager struct {
	flags     map[string]*Flag
	listeners map[string][]ChangeListener
	mu        sync.RWMutex
	provider  ConfigProvider
}

// ChangeListener is a function that is called when a flag changes
type ChangeListener func(name string, oldValue, newValue interface{})

// ConfigProvider provides flag configuration from external sources
type ConfigProvider interface {
	GetFlag(name string) (*Flag, error)
	GetAllFlags() (map[string]*Flag, error)
}

// NewFeatureFlagManager creates a new FeatureFlagManager
func NewFeatureFlagManager() *FeatureFlagManager {
	return &FeatureFlagManager{
		flags:     make(map[string]*Flag),
		listeners: make(map[string][]ChangeListener),
	}
}

// NewFeatureFlagManagerWithProvider creates a manager with a config provider
func NewFeatureFlagManagerWithProvider(provider ConfigProvider) *FeatureFlagManager {
	m := NewFeatureFlagManager()
	m.provider = provider
	return m
}

// RegisterBooleanFlag registers a new boolean flag
func (m *FeatureFlagManager) RegisterBooleanFlag(name string, defaultValue bool, description string) error {
	return m.registerFlag(FlagConfig{
		Name:        name,
		Type:        FlagTypeBoolean,
		Default:     defaultValue,
		Description: description,
		Enabled:     true,
	})
}

// RegisterStringFlag registers a new string flag
func (m *FeatureFlagManager) RegisterStringFlag(name string, defaultValue string, description string) error {
	return m.registerFlag(FlagConfig{
		Name:        name,
		Type:        FlagTypeString,
		Default:     defaultValue,
		Description: description,
		Enabled:     true,
	})
}

// RegisterIntFlag registers a new integer flag
func (m *FeatureFlagManager) RegisterIntFlag(name string, defaultValue int, description string) error {
	return m.registerFlag(FlagConfig{
		Name:        name,
		Type:        FlagTypeInt,
		Default:     defaultValue,
		Description: description,
		Enabled:     true,
	})
}

// RegisterFloatFlag registers a new float flag
func (m *FeatureFlagManager) RegisterFloatFlag(name string, defaultValue float64, description string) error {
	return m.registerFlag(FlagConfig{
		Name:        name,
		Type:        FlagTypeFloat,
		Default:     defaultValue,
		Description: description,
		Enabled:     true,
	})
}

// RegisterJSONFlag registers a new JSON flag
func (m *FeatureFlagManager) RegisterJSONFlag(name string, defaultValue interface{}, description string) error {
	return m.registerFlag(FlagConfig{
		Name:        name,
		Type:        FlagTypeJSON,
		Default:     defaultValue,
		Description: description,
		Enabled:     true,
	})
}

// registerFlag registers a new flag
func (m *FeatureFlagManager) registerFlag(config FlagConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.flags[config.Name]; exists {
		return fmt.Errorf("flag %q already registered", config.Name)
	}

	flag := &Flag{
		Name:        config.Name,
		Type:        config.Type,
		Value:       config.Default,
		Default:     config.Default,
		Description: config.Description,
		Enabled:     config.Enabled,
		UpdatedAt:   time.Now(),
	}

	m.flags[config.Name] = flag
	return nil
}

// GetBool returns the boolean value of a flag
func (m *FeatureFlagManager) GetBool(name string) bool {
	flag := m.getFlag(name)
	if flag == nil {
		return false
	}
	flag.mu.RLock()
	defer flag.mu.RUnlock()

	if val, ok := flag.Value.(bool); ok {
		return val
	}
	if val, ok := flag.Default.(bool); ok {
		return val
	}
	return false
}

// GetString returns the string value of a flag
func (m *FeatureFlagManager) GetString(name string) string {
	flag := m.getFlag(name)
	if flag == nil {
		return ""
	}
	flag.mu.RLock()
	defer flag.mu.RUnlock()

	if val, ok := flag.Value.(string); ok {
		return val
	}
	if val, ok := flag.Default.(string); ok {
		return val
	}
	return ""
}

// GetInt returns the integer value of a flag
func (m *FeatureFlagManager) GetInt(name string) int {
	flag := m.getFlag(name)
	if flag == nil {
		return 0
	}
	flag.mu.RLock()
	defer flag.mu.RUnlock()

	switch v := flag.Value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}

	switch v := flag.Default.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}

// GetFloat returns the float value of a flag
func (m *FeatureFlagManager) GetFloat(name string) float64 {
	flag := m.getFlag(name)
	if flag == nil {
		return 0.0
	}
	flag.mu.RLock()
	defer flag.mu.RUnlock()

	switch v := flag.Value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}

	switch v := flag.Default.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}
	return 0.0
}

// GetJSON returns the JSON value of a flag
func (m *FeatureFlagManager) GetJSON(name string) interface{} {
	flag := m.getFlag(name)
	if flag == nil {
		return nil
	}
	flag.mu.RLock()
	defer flag.mu.RUnlock()

	if flag.Value != nil {
		return flag.Value
	}
	return flag.Default
}

// IsEnabled returns whether a flag is enabled
func (m *FeatureFlagManager) IsEnabled(name string) bool {
	flag := m.getFlag(name)
	if flag == nil {
		return false
	}
	flag.mu.RLock()
	defer flag.mu.RUnlock()
	return flag.Enabled
}

// SetBool sets the boolean value of a flag
func (m *FeatureFlagManager) SetBool(name string, value bool) error {
	return m.setValue(name, value)
}

// SetString sets the string value of a flag
func (m *FeatureFlagManager) SetString(name string, value string) error {
	return m.setValue(name, value)
}

// SetInt sets the integer value of a flag
func (m *FeatureFlagManager) SetInt(name string, value int) error {
	return m.setValue(name, value)
}

// SetFloat sets the float value of a flag
func (m *FeatureFlagManager) SetFloat(name string, value float64) error {
	return m.setValue(name, value)
}

// SetJSON sets the JSON value of a flag
func (m *FeatureFlagManager) SetJSON(name string, value interface{}) error {
	return m.setValue(name, value)
}

// SetEnabled sets whether a flag is enabled
func (m *FeatureFlagManager) SetEnabled(name string, enabled bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	flag, exists := m.flags[name]
	if !exists {
		return fmt.Errorf("flag %q not found", name)
	}

	flag.mu.Lock()
	defer flag.mu.Unlock()
	flag.Enabled = enabled
	flag.UpdatedAt = time.Now()

	return nil
}

// setValue sets the value of a flag and notifies listeners
func (m *FeatureFlagManager) setValue(name string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	flag, exists := m.flags[name]
	if !exists {
		return fmt.Errorf("flag %q not found", name)
	}

	flag.mu.Lock()
	defer flag.mu.Unlock()

	oldValue := flag.Value
	flag.Value = value
	flag.UpdatedAt = time.Now()

	m.notifyListeners(name, oldValue, value)

	return nil
}

// getFlag returns a flag by name
func (m *FeatureFlagManager) getFlag(name string) *Flag {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.flags[name]
}

// RegisterChangeListener registers a listener for flag changes
func (m *FeatureFlagManager) RegisterChangeListener(name string, listener ChangeListener) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.listeners[name] == nil {
		m.listeners[name] = make([]ChangeListener, 0)
	}
	m.listeners[name] = append(m.listeners[name], listener)
}

// notifyListeners notifies all listeners of a flag change
func (m *FeatureFlagManager) notifyListeners(name string, oldValue, newValue interface{}) {
	m.mu.RLock()
	listeners := m.listeners[name]
	m.mu.RUnlock()

	for _, listener := range listeners {
		listener(name, oldValue, newValue)
	}
}

// GetFlag returns a flag by name
func (m *FeatureFlagManager) GetFlag(name string) (*Flag, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	flag, exists := m.flags[name]
	if !exists {
		return nil, fmt.Errorf("flag %q not found", name)
	}

	flag.mu.RLock()
	defer flag.mu.RUnlock()

	return &Flag{
		Name:        flag.Name,
		Type:        flag.Type,
		Value:       flag.Value,
		Default:     flag.Default,
		Description: flag.Description,
		Enabled:     flag.Enabled,
		UpdatedAt:   flag.UpdatedAt,
	}, nil
}

// GetAllFlags returns all registered flags
func (m *FeatureFlagManager) GetAllFlags() map[string]*Flag {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*Flag, len(m.flags))
	for name, flag := range m.flags {
		flag.mu.RLock()
		result[name] = &Flag{
			Name:        flag.Name,
			Type:        flag.Type,
			Value:       flag.Value,
			Default:     flag.Default,
			Description: flag.Description,
			Enabled:     flag.Enabled,
			UpdatedAt:   flag.UpdatedAt,
		}
		flag.mu.RUnlock()
	}
	return result
}

// LoadFromEnvironment loads flag values from environment variables
func (m *FeatureFlagManager) LoadFromEnvironment(prefix string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	prefix = strings.ToUpper(prefix)
	for name, flag := range m.flags {
		envName := prefix + "_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))

		envValue := os.Getenv(envName)
		if envValue == "" {
			continue
		}

		var parsedValue interface{}
		var err error

		switch flag.Type {
		case FlagTypeBoolean:
			parsedValue, err = strconv.ParseBool(envValue)
		case FlagTypeInt:
			var intVal int64
			intVal, err = strconv.ParseInt(envValue, 10, 64)
			parsedValue = int(intVal)
		case FlagTypeFloat:
			parsedValue, err = strconv.ParseFloat(envValue, 64)
		case FlagTypeString:
			parsedValue = envValue
		case FlagTypeJSON:
			err = json.Unmarshal([]byte(envValue), &parsedValue)
		}

		if err != nil {
			return fmt.Errorf("failed to parse env var %s: %w", envName, err)
		}

		flag.Value = parsedValue
		flag.UpdatedAt = time.Now()
	}

	return nil
}

// LoadFromFile loads flag configuration from a JSON file
func (m *FeatureFlagManager) LoadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var fileFlags map[string]interface{}
	if err := json.Unmarshal(data, &fileFlags); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for name, value := range fileFlags {
		flag, exists := m.flags[name]
		if !exists {
			continue
		}

		flag.Value = value
		flag.UpdatedAt = time.Now()
	}

	return nil
}

// SaveToFile saves flag values to a JSON file
func (m *FeatureFlagManager) SaveToFile(filepath string) error {
	flags := m.GetAllFlags()

	data := make(map[string]interface{})
	for name, flag := range flags {
		data[name] = flag.Value
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// EvaluateWithContext evaluates a flag with additional context
func (m *FeatureFlagManager) EvaluateWithContext(name string, context map[string]interface{}) (interface{}, error) {
	flag, err := m.GetFlag(name)
	if err != nil {
		return nil, err
	}

	if !flag.Enabled {
		return flag.Default, nil
	}

	switch v := flag.Value.(type) {
	case bool:
		return evaluateBoolWithContext(v, context), nil
	case string:
		return evaluateStringWithContext(v, context), nil
	case int, int64, float64:
		return v, nil
	default:
		return v, nil
	}
}

// evaluateBoolWithContext evaluates a boolean with context
func evaluateBoolWithContext(value bool, context map[string]interface{}) bool {
	if !value {
		return false
	}

	if context == nil {
		return value
	}

	return value
}

// evaluateStringWithContext evaluates a string with context
func evaluateStringWithContext(value string, context map[string]interface{}) string {
	if context == nil {
		return value
	}

	result := value
	for k, v := range context {
		placeholder := fmt.Sprintf("{%s}", k)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", v))
	}

	return result
}

// UnmarshalFlagValue unmarshals a flag value to a specific type
func UnmarshalFlagValue(flag *Flag, target interface{}) error {
	if flag == nil {
		return fmt.Errorf("flag is nil")
	}

	data, err := json.Marshal(flag.Value)
	if err != nil {
		return fmt.Errorf("failed to marshal flag value: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal to target type: %w", err)
	}

	return nil
}

// ValidateFlagType validates that a value matches the expected flag type
func ValidateFlagType(value interface{}, expectedType FlagType) bool {
	if value == nil {
		return false
	}

	switch expectedType {
	case FlagTypeBoolean:
		_, ok := value.(bool)
		return ok
	case FlagTypeString:
		_, ok := value.(string)
		return ok
	case FlagTypeInt:
		switch value.(type) {
		case int, int64:
			return true
		case float64:
			return float64(int(value.(float64))) == value.(float64)
		}
		return false
	case FlagTypeFloat:
		_, ok := value.(float64)
		if !ok {
			_, ok = value.(int)
		}
		return ok
	case FlagTypeJSON:
		return true
	default:
		return false
	}
}

// GetFlagType returns the FlagType for a given Go type
func GetFlagType(v interface{}) FlagType {
	if v == nil {
		return FlagTypeJSON
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Bool:
		return FlagTypeBoolean
	case reflect.String:
		return FlagTypeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return FlagTypeInt
	case reflect.Float32, reflect.Float64:
		return FlagTypeFloat
	default:
		return FlagTypeJSON
	}
}

// ToJSON converts flags to JSON
func (m *FeatureFlagManager) ToJSON() ([]byte, error) {
	flags := m.GetAllFlags()
	return json.MarshalIndent(flags, "", "  ")
}

// FromJSON loads flags from JSON
func FromJSON(data []byte) (map[string]*Flag, error) {
	var flags map[string]*Flag
	if err := json.Unmarshal(data, &flags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal flags: %w", err)
	}
	return flags, nil
}
