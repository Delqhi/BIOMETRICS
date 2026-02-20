package envconfig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotAPointer      = errors.New("config must be a pointer")
	ErrNotAStruct       = errors.New("config must be a struct")
	ErrUnsupportedType  = errors.New("unsupported field type")
	ErrRequiredField    = errors.New("required field is empty")
	ErrInvalidValue     = errors.New("invalid value for field")
	ErrValidationFailed = errors.New("validation failed")
)

type fieldInfo struct {
	name     string
	envVar   string
	defValue string
	required bool
	desc     string
}

type Config struct {
	Prefix      string
	Strict      bool
	AllowUnused bool
}

type Option func(*Config)

func WithPrefix(prefix string) Option {
	return func(c *Config) {
		c.Prefix = prefix
	}
}

func WithStrict(strict bool) Option {
	return func(c *Config) {
		c.Strict = strict
	}
}

func WithAllowUnused(allow bool) Option {
	return func(c *Config) {
		c.AllowUnused = allow
	}
}

func Process(prefix string, spec interface{}) error {
	return ProcessWithOptions(spec, WithPrefix(prefix))
}

func ProcessWithOptions(spec interface{}, opts ...Option) error {
	config := &Config{
		Prefix:      "",
		Strict:      false,
		AllowUnused: true,
	}

	for _, opt := range opts {
		opt(config)
	}

	return process(config, spec)
}

func process(config *Config, spec interface{}) error {
	val := reflect.ValueOf(spec)
	if val.Kind() != reflect.Ptr {
		return ErrNotAPointer
	}

	if val.IsNil() {
		return ErrNotAPointer
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !field.CanSet() {
			continue
		}

		info := parseFieldTag(fieldType)

		envVar := info.envVar
		if envVar == "" {
			envVar = fieldToEnvVar(fieldType.Name, config.Prefix)
		}

		envValue := os.Getenv(envVar)

		if envValue == "" {
			if info.defValue != "" {
				envValue = info.defValue
			} else if info.required {
				return fmt.Errorf("%w: %s", ErrRequiredField, envVar)
			} else {
				continue
			}
		}

		if err := setFieldValue(field, envValue, envVar); err != nil {
			return fmt.Errorf("field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

func parseFieldTag(field reflect.StructField) *fieldInfo {
	info := &fieldInfo{
		name: field.Name,
	}

	tag := field.Tag.Get("env")
	if tag == "" {
		return info
	}

	parts := strings.Split(tag, ",")
	if len(parts) > 0 && parts[0] != "" {
		info.envVar = parts[0]
	}

	for _, part := range parts[1:] {
		kv := strings.SplitN(part, "=", 2)
		switch kv[0] {
		case "required":
			info.required = true
		case "default":
			if len(kv) > 1 {
				info.defValue = kv[1]
			}
		case "desc":
			if len(kv) > 1 {
				info.desc = kv[1]
			}
		}
	}

	defTag := field.Tag.Get("default")
	if defTag != "" && info.defValue == "" {
		info.defValue = defTag
	}

	return info
}

func setFieldValue(field reflect.Value, value, envVar string) error {
	if !field.CanSet() {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		v, err := parseBool(value)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidValue, err)
		}
		field.SetBool(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidValue, err)
			}
			field.SetInt(int64(d))
		} else {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidValue, err)
			}
			field.SetInt(v)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidValue, err)
		}
		field.SetUint(v)

	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidValue, err)
		}
		field.SetFloat(v)

	case reflect.Slice:
		return setSliceValue(field, value)

	case reflect.Map:
		return setMapValue(field, value)

	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setFieldValue(field.Elem(), value, envVar)

	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedType, field.Kind())
	}

	return nil
}

func setSliceValue(field reflect.Value, value string) error {
	items := strings.Split(value, ",")
	slice := reflect.MakeSlice(field.Type(), len(items), len(items))

	elemType := field.Type().Elem()

	for i, item := range items {
		item = strings.TrimSpace(item)
		elem := slice.Index(i)

		switch elemType.Kind() {
		case reflect.String:
			elem.SetString(item)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(item, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: slice item %d: %s", ErrInvalidValue, i, err)
			}
			elem.SetInt(v)
		case reflect.Bool:
			v, err := parseBool(item)
			if err != nil {
				return fmt.Errorf("%w: slice item %d: %s", ErrInvalidValue, i, err)
			}
			elem.SetBool(v)
		default:
			return fmt.Errorf("%w: slice element type %s", ErrUnsupportedType, elemType.Kind())
		}
	}

	field.Set(slice)
	return nil
}

func setMapValue(field reflect.Value, value string) error {
	pairs := strings.Split(value, ",")
	m := reflect.MakeMap(field.Type())

	keyType := field.Type().Key()
	elemType := field.Type().Elem()

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			return fmt.Errorf("%w: invalid map entry format, expected key:value", ErrInvalidValue)
		}

		key := reflect.ValueOf(kv[0])
		if keyType != reflect.TypeOf("") {
			key = reflect.ValueOf(kv[0]).Convert(keyType)
		}

		var elem reflect.Value
		switch elemType.Kind() {
		case reflect.String:
			elem = reflect.ValueOf(kv[1])
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				return fmt.Errorf("%w: map value: %s", ErrInvalidValue, err)
			}
			elem = reflect.ValueOf(v).Convert(elemType)
		default:
			return fmt.Errorf("%w: map value type %s", ErrUnsupportedType, elemType.Kind())
		}

		m.SetMapIndex(key, elem)
	}

	field.Set(m)
	return nil
}

func parseBool(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "on", "enabled":
		return true, nil
	case "false", "0", "no", "off", "disabled", "":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", value)
	}
}

func fieldToEnvVar(name, prefix string) string {
	var result []rune
	for i, r := range name {
		if i > 0 && isUpper(r) && (i < len(name)-1 && isLower(rune(name[i+1])) || i > 0 && isLower(rune(name[i-1]))) {
			result = append(result, '_')
		}
		result = append(result, toUpper(r))
	}

	envVar := string(result)
	if prefix != "" {
		envVar = strings.ToUpper(prefix) + "_" + envVar
	}

	return envVar
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	return r
}

func MustProcess(prefix string, spec interface{}) {
	if err := Process(prefix, spec); err != nil {
		panic(err)
	}
}

func Getenv(key, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defValue
}

func GetenvInt(key string, defValue int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return defValue
}

func GetenvBool(key string, defValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if v, err := parseBool(value); err == nil {
			return v
		}
	}
	return defValue
}

func GetenvDuration(key string, defValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defValue
}

func GetenvSlice(key string, defValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defValue
}

func Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func Unsetenv(key string) error {
	return os.Unsetenv(key)
}

type EnvVar struct {
	Key      string
	Value    string
	HasValue bool
}

func List(prefix string) []EnvVar {
	var result []EnvVar

	for _, pair := range os.Environ() {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			if prefix == "" || strings.HasPrefix(kv[0], prefix) {
				result = append(result, EnvVar{
					Key:      kv[0],
					Value:    kv[1],
					HasValue: true,
				})
			}
		}
	}

	return result
}

func ValidateRequired(spec interface{}) error {
	val := reflect.ValueOf(spec)
	if val.Kind() != reflect.Ptr {
		return ErrNotAPointer
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("env")
		if strings.Contains(tag, ",required") {
			if isZero(field) {
				envVar := strings.Split(tag, ",")[0]
				if envVar == "" {
					envVar = fieldType.Name
				}
				return fmt.Errorf("%w: %s", ErrRequiredField, envVar)
			}
		}
	}

	return nil
}

func isZero(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func Export(spec interface{}, prefix string) map[string]string {
	result := make(map[string]string)

	val := reflect.ValueOf(spec)
	if val.Kind() != reflect.Ptr {
		return result
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return result
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !field.CanInterface() {
			continue
		}

		envVar := fieldToEnvVar(fieldType.Name, prefix)
		tag := fieldType.Tag.Get("env")
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				envVar = parts[0]
			}
		}

		result[envVar] = fmt.Sprintf("%v", field.Interface())
	}

	return result
}
