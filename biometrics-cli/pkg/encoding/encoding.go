package encoding

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrDecodingFailed   = errors.New("decoding failed")
	ErrEncodingFailed   = errors.New("encoding failed")
	ErrUnsupportedType  = errors.New("unsupported type")
	ErrCompressionError = errors.New("compression error")
)

type EncodingType string

const (
	EncodingBase64       EncodingType = "base64"
	EncodingBase64URL    EncodingType = "base64url"
	EncodingBase64Raw    EncodingType = "base64raw"
	EncodingBase64RawURL EncodingType = "base64rawurl"
	EncodingHex          EncodingType = "hex"
	EncodingJSON         EncodingType = "json"
	EncodingGzip         EncodingType = "gzip"
	EncodingGzipBase64   EncodingType = "gzipbase64"
)

type Encoder struct {
	encodingType EncodingType
}

func NewEncoder(encodingType EncodingType) *Encoder {
	return &Encoder{encodingType: encodingType}
}

func (e *Encoder) Encode(data []byte) (string, error) {
	switch e.encodingType {
	case EncodingBase64:
		return base64.StdEncoding.EncodeToString(data), nil
	case EncodingBase64URL:
		return base64.URLEncoding.EncodeToString(data), nil
	case EncodingBase64Raw:
		return base64.RawStdEncoding.EncodeToString(data), nil
	case EncodingBase64RawURL:
		return base64.RawURLEncoding.EncodeToString(data), nil
	case EncodingHex:
		return hex.EncodeToString(data), nil
	case EncodingJSON:
		b, err := json.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrEncodingFailed, err)
		}
		return string(b), nil
	case EncodingGzip:
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)
		if _, err := w.Write(data); err != nil {
			return "", fmt.Errorf("%w: %v", ErrCompressionError, err)
		}
		if err := w.Close(); err != nil {
			return "", fmt.Errorf("%w: %v", ErrCompressionError, err)
		}
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	case EncodingGzipBase64:
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)
		if _, err := w.Write(data); err != nil {
			return "", fmt.Errorf("%w: %v", ErrCompressionError, err)
		}
		if err := w.Close(); err != nil {
			return "", fmt.Errorf("%w: %v", ErrCompressionError, err)
		}
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	default:
		return "", ErrUnsupportedType
	}
}

func (e *Encoder) Decode(data string) ([]byte, error) {
	switch e.encodingType {
	case EncodingBase64:
		return base64.StdEncoding.DecodeString(data)
	case EncodingBase64URL:
		return base64.URLEncoding.DecodeString(data)
	case EncodingBase64Raw:
		return base64.RawStdEncoding.DecodeString(data)
	case EncodingBase64RawURL:
		return base64.RawURLEncoding.DecodeString(data)
	case EncodingHex:
		return hex.DecodeString(data)
	case EncodingJSON:
		var result []byte
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDecodingFailed, err)
		}
		return result, nil
	case EncodingGzip, EncodingGzipBase64:
		compressed, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDecodingFailed, err)
		}
		r, err := gzip.NewReader(bytes.NewReader(compressed))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDecodingFailed, err)
		}
		defer r.Close()
		result, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDecodingFailed, err)
		}
		return result, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func Base64URLEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func Base64URLDecode(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}

func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

func HexDecode(data string) ([]byte, error) {
	return hex.DecodeString(data)
}

func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCompressionError, err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCompressionError, err)
	}
	return buf.Bytes(), nil
}

func GzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecodingFailed, err)
	}
	defer r.Close()
	result, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecodingFailed, err)
	}
	return result, nil
}

func GzipCompressBase64(data []byte) (string, error) {
	compressed, err := GzipCompress(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(compressed), nil
}

func GzipDecompressBase64(data string) ([]byte, error) {
	compressed, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecodingFailed, err)
	}
	return GzipDecompress(compressed)
}

func JSONEncode(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrEncodingFailed, err)
	}
	return string(data), nil
}

func JSONEncodeIndent(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrEncodingFailed, err)
	}
	return string(data), nil
}

func JSONDecode(data string, v interface{}) error {
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return fmt.Errorf("%w: %v", ErrDecodingFailed, err)
	}
	return nil
}

func MustEncode(data []byte, encodingType EncodingType) string {
	encoder := NewEncoder(encodingType)
	result, err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	return result
}

func MustDecode(data string, encodingType EncodingType) []byte {
	encoder := NewEncoder(encodingType)
	result, err := encoder.Decode(data)
	if err != nil {
		panic(err)
	}
	return result
}

func DetectEncoding(data string) EncodingType {
	if _, err := base64.StdEncoding.DecodeString(data); err == nil {
		return EncodingBase64
	}
	if _, err := base64.URLEncoding.DecodeString(data); err == nil {
		return EncodingBase64URL
	}
	if _, err := hex.DecodeString(data); err == nil && len(data)%2 == 0 {
		return EncodingHex
	}
	if json.Valid([]byte(data)) {
		return EncodingJSON
	}
	return ""
}

func AutoDecode(data string) ([]byte, EncodingType, error) {
	encodingType := DetectEncoding(data)
	if encodingType == "" {
		return nil, "", ErrUnsupportedType
	}

	encoder := NewEncoder(encodingType)
	result, err := encoder.Decode(data)
	if err != nil {
		return nil, encodingType, err
	}
	return result, encodingType, nil
}

func ValidateBase64(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}

func ValidateHex(data string) bool {
	_, err := hex.DecodeString(data)
	return err == nil && len(data)%2 == 0
}

func ValidateJSON(data string) bool {
	return json.Valid([]byte(data))
}

func StripPadding(data string) string {
	return strings.TrimRight(data, "=")
}

func AddPadding(data string) string {
	if mod := len(data) % 4; mod != 0 {
		data += strings.Repeat("=", 4-mod)
	}
	return data
}

func BytesToString(data []byte) string {
	return string(data)
}

func StringToBytes(data string) []byte {
	return []byte(data)
}

func EncodeMap(m map[string]interface{}, encodingType EncodingType) (map[string]string, error) {
	result := make(map[string]string)
	encoder := NewEncoder(encodingType)

	for k, v := range m {
		var data []byte
		var err error

		switch val := v.(type) {
		case string:
			data = []byte(val)
		case []byte:
			data = val
		default:
			data, err = json.Marshal(val)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal value for key %s: %w", k, err)
			}
		}

		encoded, err := encoder.Encode(data)
		if err != nil {
			return nil, fmt.Errorf("failed to encode value for key %s: %w", k, err)
		}
		result[k] = encoded
	}

	return result, nil
}

func DecodeMap(m map[string]string, encodingType EncodingType) (map[string][]byte, error) {
	result := make(map[string][]byte)
	encoder := NewEncoder(encodingType)

	for k, v := range m {
		decoded, err := encoder.Decode(v)
		if err != nil {
			return nil, fmt.Errorf("failed to decode value for key %s: %w", k, err)
		}
		result[k] = decoded
	}

	return result, nil
}

type MultiEncoder struct {
	encodings []EncodingType
}

func NewMultiEncoder(encodings ...EncodingType) *MultiEncoder {
	return &MultiEncoder{encodings: encodings}
}

func (m *MultiEncoder) Encode(data []byte) (string, error) {
	var result = data
	var err error

	for i := len(m.encodings) - 1; i >= 0; i-- {
		encoder := NewEncoder(m.encodings[i])
		encoded, err := encoder.Encode(result)
		if err != nil {
			return "", fmt.Errorf("encoding failed at step %d: %w", i, err)
		}
		result = []byte(encoded)
	}

	return string(result), err
}

func (m *MultiEncoder) Decode(data string) ([]byte, error) {
	var result = []byte(data)

	for _, encoding := range m.encodings {
		encoder := NewEncoder(encoding)
		decoded, err := encoder.Decode(string(result))
		if err != nil {
			return nil, fmt.Errorf("decoding failed for %s: %w", encoding, err)
		}
		result = decoded
	}

	return result, nil
}

func IsPrintable(data []byte) bool {
	for _, b := range data {
		if b < 32 || b > 126 {
			if b != '\n' && b != '\r' && b != '\t' {
				return false
			}
		}
	}
	return true
}

func SafeString(data []byte) string {
	if IsPrintable(data) {
		return string(data)
	}
	return HexEncode(data)
}
