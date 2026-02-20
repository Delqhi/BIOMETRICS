package envconfig

import (
	"os"
	"reflect"
	"testing"
	"time"
)

type TestConfig struct {
	StringVar  string        `env:"TEST_STRING" default:"default"`
	IntVar     int           `env:"TEST_INT" default:"42"`
	BoolVar    bool          `env:"TEST_BOOL" default:"true"`
	FloatVar   float64       `env:"TEST_FLOAT" default:"3.14"`
	Duration   time.Duration `env:"TEST_DURATION" default:"5s"`
	StringList []string      `env:"TEST_LIST" default:"a,b,c"`
	Required   string        `env:"TEST_REQUIRED,required"`
	NoTag      string
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected TestConfig
		hasError bool
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: TestConfig{
				StringVar:  "default",
				IntVar:     42,
				BoolVar:    true,
				FloatVar:   3.14,
				Duration:   5 * time.Second,
				StringList: []string{"a", "b", "c"},
			},
			hasError: true,
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"TEST_STRING":   "custom",
				"TEST_INT":      "100",
				"TEST_BOOL":     "false",
				"TEST_FLOAT":    "2.71",
				"TEST_DURATION": "10m",
				"TEST_LIST":     "x,y,z",
				"TEST_REQUIRED": "present",
			},
			expected: TestConfig{
				StringVar:  "custom",
				IntVar:     100,
				BoolVar:    false,
				FloatVar:   2.71,
				Duration:   10 * time.Minute,
				StringList: []string{"x", "y", "z"},
				Required:   "present",
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			var config TestConfig
			err := Process("", &config)

			if tt.hasError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Process() error: %v", err)
			}

			if config.StringVar != tt.expected.StringVar {
				t.Errorf("StringVar = %q, want %q", config.StringVar, tt.expected.StringVar)
			}
			if config.IntVar != tt.expected.IntVar {
				t.Errorf("IntVar = %d, want %d", config.IntVar, tt.expected.IntVar)
			}
			if config.BoolVar != tt.expected.BoolVar {
				t.Errorf("BoolVar = %v, want %v", config.BoolVar, tt.expected.BoolVar)
			}
			if config.FloatVar != tt.expected.FloatVar {
				t.Errorf("FloatVar = %v, want %v", config.FloatVar, tt.expected.FloatVar)
			}
			if config.Duration != tt.expected.Duration {
				t.Errorf("Duration = %v, want %v", config.Duration, tt.expected.Duration)
			}
		})
	}
}

func TestProcessWithPrefix(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_TEST_VAR", "value")

	type Config struct {
		TestVar string
	}

	var config Config
	err := Process("APP", &config)
	if err != nil {
		t.Fatalf("Process() error: %v", err)
	}

	if config.TestVar != "value" {
		t.Errorf("TestVar = %q, want %q", config.TestVar, "value")
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		hasError bool
	}{
		{"true", true, false},
		{"True", true, false},
		{"TRUE", true, false},
		{"1", true, false},
		{"yes", true, false},
		{"on", true, false},
		{"enabled", true, false},
		{"false", false, false},
		{"False", false, false},
		{"0", false, false},
		{"no", false, false},
		{"off", false, false},
		{"invalid", false, true},
	}

	for _, tt := range tests {
		result, err := parseBool(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("parseBool(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("parseBool(%q) error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("parseBool(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		}
	}
}

func TestSetSliceValue(t *testing.T) {
	type Config struct {
		Strings []string
		Ints    []int
		Bools   []bool
	}

	tests := []struct {
		name     string
		input    string
		field    string
		expected interface{}
	}{
		{"strings", "a,b,c", "Strings", []string{"a", "b", "c"}},
		{"ints", "1,2,3", "Ints", []int{1, 2, 3}},
		{"bools", "true,false,1", "Bools", []bool{true, false, true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config Config
			var field reflect.Value

			switch tt.field {
			case "Strings":
				field = reflect.ValueOf(&config.Strings).Elem()
			case "Ints":
				field = reflect.ValueOf(&config.Ints).Elem()
			case "Bools":
				field = reflect.ValueOf(&config.Bools).Elem()
			}

			if err := setSliceValue(field, tt.input); err != nil {
				t.Fatalf("setSliceValue() error: %v", err)
			}

			if !reflect.DeepEqual(field.Interface(), tt.expected) {
				t.Errorf("got %v, want %v", field.Interface(), tt.expected)
			}
		})
	}
}

func TestGetenv(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_VAR", "value")

	if v := Getenv("TEST_VAR", "default"); v != "value" {
		t.Errorf("Getenv() = %q, want %q", v, "value")
	}

	if v := Getenv("NONEXISTENT", "default"); v != "default" {
		t.Errorf("Getenv() = %q, want %q", v, "default")
	}
}

func TestGetenvInt(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_INT", "42")

	if v := GetenvInt("TEST_INT", 0); v != 42 {
		t.Errorf("GetenvInt() = %d, want %d", v, 42)
	}

	if v := GetenvInt("NONEXISTENT", 100); v != 100 {
		t.Errorf("GetenvInt() = %d, want %d", v, 100)
	}
}

func TestGetenvBool(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_BOOL", "true")

	if v := GetenvBool("TEST_BOOL", false); v != true {
		t.Errorf("GetenvBool() = %v, want %v", v, true)
	}

	if v := GetenvBool("NONEXISTENT", true); v != true {
		t.Errorf("GetenvBool() = %v, want %v", v, true)
	}
}

func TestGetenvDuration(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_DURATION", "5m")

	if v := GetenvDuration("TEST_DURATION", 0); v != 5*time.Minute {
		t.Errorf("GetenvDuration() = %v, want %v", v, 5*time.Minute)
	}

	if v := GetenvDuration("NONEXISTENT", time.Hour); v != time.Hour {
		t.Errorf("GetenvDuration() = %v, want %v", v, time.Hour)
	}
}

func TestGetenvSlice(t *testing.T) {
	os.Clearenv()
	os.Setenv("TEST_SLICE", "a,b,c")

	result := GetenvSlice("TEST_SLICE", nil)
	if len(result) != 3 {
		t.Errorf("GetenvSlice() length = %d, want 3", len(result))
	}

	result = GetenvSlice("NONEXISTENT", []string{"default"})
	if len(result) != 1 || result[0] != "default" {
		t.Errorf("GetenvSlice() = %v, want [default]", result)
	}
}

func TestList(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_VAR1", "value1")
	os.Setenv("APP_VAR2", "value2")
	os.Setenv("OTHER_VAR", "other")

	vars := List("APP")
	if len(vars) != 2 {
		t.Errorf("List() returned %d vars, want 2", len(vars))
	}

	vars = List("")
	if len(vars) < 3 {
		t.Errorf("List() returned %d vars, want at least 3", len(vars))
	}
}

func TestNotAPointer(t *testing.T) {
	var config TestConfig
	err := Process("", config)
	if err != ErrNotAPointer {
		t.Errorf("expected ErrNotAPointer, got %v", err)
	}
}

func TestNotAStruct(t *testing.T) {
	var config string
	err := Process("", &config)
	if err != ErrNotAStruct {
		t.Errorf("expected ErrNotAStruct, got %v", err)
	}
}

func TestMustProcess(t *testing.T) {
	os.Clearenv()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustProcess should panic on error")
		}
	}()

	type Config struct {
		Required string `env:"REQ,required"`
	}

	var config Config
	MustProcess("", &config)
}

func TestExport(t *testing.T) {
	type Config struct {
		Name  string
		Value int
	}

	config := Config{
		Name:  "test",
		Value: 42,
	}

	exported := Export(&config, "APP")

	if exported["APP_NAME"] != "test" {
		t.Errorf("APP_NAME = %q, want %q", exported["APP_NAME"], "test")
	}

	if exported["APP_VALUE"] != "42" {
		t.Errorf("APP_VALUE = %q, want %q", exported["APP_VALUE"], "42")
	}
}

func TestFieldToEnvVar(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		expected string
	}{
		{"CamelCase", "", "CAMEL_CASE"},
		{"Test", "APP", "APP_TEST"},
	}

	for _, tt := range tests {
		result := fieldToEnvVar(tt.name, tt.prefix)
		if result != tt.expected {
			t.Errorf("fieldToEnvVar(%q, %q) = %q, want %q",
				tt.name, tt.prefix, result, tt.expected)
		}
	}
}
