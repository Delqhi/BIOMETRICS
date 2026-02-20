package encoding

import (
	"bytes"
	"compress/gzip"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte("hello"), "aGVsbG8="},
		{[]byte("world"), "d29ybGQ="},
		{[]byte(""), ""},
		{[]byte("test data"), "dGVzdCBkYXRh"},
	}

	for _, tt := range tests {
		result := Base64Encode(tt.input)
		if result != tt.expected {
			t.Errorf("Base64Encode(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestBase64Decode(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
		hasError bool
	}{
		{"aGVsbG8=", []byte("hello"), false},
		{"d29ybGQ=", []byte("world"), false},
		{"", []byte{}, false},
		{"invalid!@", nil, true},
	}

	for _, tt := range tests {
		result, err := Base64Decode(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("Base64Decode(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("Base64Decode(%q) unexpected error: %v", tt.input, err)
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Base64Decode(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		}
	}
}

func TestHexEncode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte{0x00, 0x01, 0x02}, "000102"},
		{[]byte{0xff, 0xfe}, "fffe"},
		{[]byte("AB"), "4142"},
	}

	for _, tt := range tests {
		result := HexEncode(tt.input)
		if result != tt.expected {
			t.Errorf("HexEncode(%v) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestHexDecode(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
		hasError bool
	}{
		{"000102", []byte{0x00, 0x01, 0x02}, false},
		{"fffe", []byte{0xff, 0xfe}, false},
		{"4142", []byte("AB"), false},
		{"invalid", nil, true},
		{"0", nil, true},
	}

	for _, tt := range tests {
		result, err := HexDecode(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("HexDecode(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("HexDecode(%q) unexpected error: %v", tt.input, err)
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("HexDecode(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		}
	}
}

func TestGzipCompressDecompress(t *testing.T) {
	tests := []struct {
		input []byte
	}{
		{[]byte("hello world")},
		{[]byte("")},
		{make([]byte, 1000)},
		{bytes.Repeat([]byte("test"), 100)},
	}

	for _, tt := range tests {
		compressed, err := GzipCompress(tt.input)
		if err != nil {
			t.Errorf("GzipCompress() error: %v", err)
			continue
		}

		decompressed, err := GzipDecompress(compressed)
		if err != nil {
			t.Errorf("GzipDecompress() error: %v", err)
			continue
		}

		if !bytes.Equal(decompressed, tt.input) {
			t.Errorf("Gzip roundtrip failed: got %v, want %v", decompressed, tt.input)
		}
	}
}

func TestGzipCompressBase64(t *testing.T) {
	input := []byte("test data for compression")

	encoded, err := GzipCompressBase64(input)
	if err != nil {
		t.Fatalf("GzipCompressBase64() error: %v", err)
	}

	decoded, err := GzipDecompressBase64(encoded)
	if err != nil {
		t.Fatalf("GzipDecompressBase64() error: %v", err)
	}

	if !bytes.Equal(decoded, input) {
		t.Errorf("GzipBase64 roundtrip failed: got %v, want %v", decoded, input)
	}
}

func TestJSONEncode(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{map[string]string{"key": "value"}, `{"key":"value"}`},
		{[]string{"a", "b"}, `["a","b"]`},
		{"string", `"string"`},
		{123, `123`},
	}

	for _, tt := range tests {
		result, err := JSONEncode(tt.input)
		if err != nil {
			t.Errorf("JSONEncode() error: %v", err)
			continue
		}
		if result != tt.expected {
			t.Errorf("JSONEncode() = %q, want %q", result, tt.expected)
		}
	}
}

func TestJSONDecode(t *testing.T) {
	var result map[string]string
	err := JSONDecode(`{"key":"value"}`, &result)
	if err != nil {
		t.Errorf("JSONDecode() error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("JSONDecode() = %v, want {key: value}", result)
	}
}

func TestEncoder(t *testing.T) {
	t.Run("Base64", func(t *testing.T) {
		encoder := NewEncoder(EncodingBase64)
		data := []byte("test")

		encoded, err := encoder.Encode(data)
		if err != nil {
			t.Fatalf("Encode() error: %v", err)
		}

		decoded, err := encoder.Decode(encoded)
		if err != nil {
			t.Fatalf("Decode() error: %v", err)
		}

		if !bytes.Equal(decoded, data) {
			t.Errorf("Roundtrip failed: got %v, want %v", decoded, data)
		}
	})

	t.Run("Hex", func(t *testing.T) {
		encoder := NewEncoder(EncodingHex)
		data := []byte{0x01, 0x02, 0x03}

		encoded, err := encoder.Encode(data)
		if err != nil {
			t.Fatalf("Encode() error: %v", err)
		}

		decoded, err := encoder.Decode(encoded)
		if err != nil {
			t.Fatalf("Decode() error: %v", err)
		}

		if !bytes.Equal(decoded, data) {
			t.Errorf("Roundtrip failed: got %v, want %v", decoded, data)
		}
	})

	t.Run("Unsupported", func(t *testing.T) {
		encoder := NewEncoder("unsupported")
		_, err := encoder.Encode([]byte("test"))
		if err != ErrUnsupportedType {
			t.Errorf("Expected ErrUnsupportedType, got %v", err)
		}
	})
}

func TestDetectEncoding(t *testing.T) {
	tests := []struct {
		input    string
		expected EncodingType
	}{
		{"aGVsbG8=", EncodingBase64},
		{"68656c6c6f", EncodingHex},
		{`{"key":"value"}`, EncodingJSON},
		{"not encoded", ""},
	}

	for _, tt := range tests {
		result := DetectEncoding(tt.input)
		if result != tt.expected {
			t.Errorf("DetectEncoding(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestAutoDecode(t *testing.T) {
	t.Run("Base64", func(t *testing.T) {
		data := []byte("hello")
		encoded := Base64Encode(data)

		decoded, encType, err := AutoDecode(encoded)
		if err != nil {
			t.Fatalf("AutoDecode() error: %v", err)
		}
		if encType != EncodingBase64 {
			t.Errorf("Encoding type = %q, want %q", encType, EncodingBase64)
		}
		if !bytes.Equal(decoded, data) {
			t.Errorf("Decoded = %v, want %v", decoded, data)
		}
	})

	t.Run("Unsupported", func(t *testing.T) {
		_, _, err := AutoDecode("not valid encoding")
		if err != ErrUnsupportedType {
			t.Errorf("Expected ErrUnsupportedType, got %v", err)
		}
	})
}

func TestValidateBase64(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"aGVsbG8=", true},
		{"invalid!@", false},
		{"", true},
	}

	for _, tt := range tests {
		result := ValidateBase64(tt.input)
		if result != tt.expected {
			t.Errorf("ValidateBase64(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestValidateHex(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"010203", true},
		{"invalid", false},
		{"0", false},
		{"", true},
	}

	for _, tt := range tests {
		result := ValidateHex(tt.input)
		if result != tt.expected {
			t.Errorf("ValidateHex(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestStripAndAddPadding(t *testing.T) {
	tests := []struct {
		input    string
		stripped string
		padded   string
	}{
		{"aGVsbG8=", "aGVsbG8", "aGVsbG8="},
		{"YQ==", "YQ", "YQ=="},
		{"YWJjZA==", "YWJjZA", "YWJjZA=="},
	}

	for _, tt := range tests {
		stripped := StripPadding(tt.input)
		if stripped != tt.stripped {
			t.Errorf("StripPadding(%q) = %q, want %q", tt.input, stripped, tt.stripped)
		}

		padded := AddPadding(tt.stripped)
		if padded != tt.padded {
			t.Errorf("AddPadding(%q) = %q, want %q", tt.stripped, padded, tt.padded)
		}
	}
}

func TestEncodeDecodeMap(t *testing.T) {
	m := map[string]interface{}{
		"key1": "value1",
		"key2": []byte("value2"),
		"key3": 123,
	}

	encoded, err := EncodeMap(m, EncodingBase64)
	if err != nil {
		t.Fatalf("EncodeMap() error: %v", err)
	}

	decoded, err := DecodeMap(encoded, EncodingBase64)
	if err != nil {
		t.Fatalf("DecodeMap() error: %v", err)
	}

	if string(decoded["key1"]) != "value1" {
		t.Errorf("key1 = %v, want value1", decoded["key1"])
	}
	if string(decoded["key2"]) != "value2" {
		t.Errorf("key2 = %v, want value2", decoded["key2"])
	}
}

func TestMultiEncoder(t *testing.T) {
	multi := NewMultiEncoder(EncodingBase64, EncodingHex)

	data := []byte("test")

	encoded, err := multi.Encode(data)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	decoded, err := multi.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode() error: %v", err)
	}

	if !bytes.Equal(decoded, data) {
		t.Errorf("Roundtrip failed: got %v, want %v", decoded, data)
	}
}

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		input    []byte
		expected bool
	}{
		{[]byte("hello"), true},
		{[]byte("hello\nworld"), true},
		{[]byte{0x00, 0x01}, false},
		{[]byte{0x7f}, false},
	}

	for _, tt := range tests {
		result := IsPrintable(tt.input)
		if result != tt.expected {
			t.Errorf("IsPrintable(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestSafeString(t *testing.T) {
	t.Run("Printable", func(t *testing.T) {
		result := SafeString([]byte("hello"))
		if result != "hello" {
			t.Errorf("SafeString() = %q, want %q", result, "hello")
		}
	})

	t.Run("NonPrintable", func(t *testing.T) {
		result := SafeString([]byte{0x00, 0x01})
		if result != "0001" {
			t.Errorf("SafeString() = %q, want %q", result, "0001")
		}
	})
}

func TestMustEncodePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustEncode should panic for unsupported type")
		}
	}()

	MustEncode([]byte("test"), "unsupported")
}

func TestMustDecodePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustDecode should panic for invalid input")
		}
	}()

	MustDecode("invalid!@", EncodingBase64)
}

func TestGzipWithRealCompressedData(t *testing.T) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write([]byte("test data"))
	w.Close()

	decompressed, err := GzipDecompress(buf.Bytes())
	if err != nil {
		t.Fatalf("GzipDecompress() error: %v", err)
	}

	if string(decompressed) != "test data" {
		t.Errorf("Decompressed = %q, want %q", decompressed, "test data")
	}
}
