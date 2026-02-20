package encryption

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Algorithm != AlgorithmAES256GCM {
		t.Errorf("Expected algorithm %s, got %s", AlgorithmAES256GCM, config.Algorithm)
	}

	if config.KDF != KDFArgon2 {
		t.Errorf("Expected KDF %s, got %s", KDFArgon2, config.KDF)
	}

	if config.KeySize != KeySize {
		t.Errorf("Expected key size %d, got %d", KeySize, config.KeySize)
	}

	if config.NonceSize != NonceSize {
		t.Errorf("Expected nonce size %d, got %d", NonceSize, config.NonceSize)
	}
}

func TestNewManager(t *testing.T) {
	config := DefaultConfig()
	m := NewManager(config)

	if m == nil {
		t.Fatal("NewManager returned nil")
	}

	if m.config.Algorithm != config.Algorithm {
		t.Errorf("Config not set correctly")
	}
}

func TestNewManagerWithKey(t *testing.T) {
	config := DefaultConfig()
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	m, err := NewManagerWithKey(config, key, "test-key")
	if err != nil {
		t.Fatalf("NewManagerWithKey failed: %v", err)
	}

	if m == nil {
		t.Fatal("NewManagerWithKey returned nil")
	}

	if m.keyHint != "test-key" {
		t.Errorf("Expected key hint 'test-key', got '%s'", m.keyHint)
	}
}

func TestNewManagerWithKeyInvalid(t *testing.T) {
	config := DefaultConfig()
	invalidKey := make([]byte, 16)

	_, err := NewManagerWithKey(config, invalidKey, "")
	if err != ErrInvalidKey {
		t.Errorf("Expected ErrInvalidKey, got %v", err)
	}
}

func TestSetMasterKey(t *testing.T) {
	config := DefaultConfig()
	m := NewManager(config)

	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	err = m.SetMasterKey(key, "my-key")
	if err != nil {
		t.Fatalf("SetMasterKey failed: %v", err)
	}

	if m.keyHint != "my-key" {
		t.Errorf("Expected key hint 'my-key', got '%s'", m.keyHint)
	}
}

func TestSetMasterKeyInvalid(t *testing.T) {
	config := DefaultConfig()
	m := NewManager(config)

	invalidKey := make([]byte, 16)

	err := m.SetMasterKey(invalidKey, "")
	if err != ErrInvalidKey {
		t.Errorf("Expected ErrInvalidKey, got %v", err)
	}
}

func TestDeriveKeyArgon2(t *testing.T) {
	config := DefaultConfig()
	config.KDF = KDFArgon2

	m := NewManager(config)

	salt := make([]byte, SaltSize)
	key := m.DeriveKey("password123", salt)

	if len(key) != KeySize {
		t.Errorf("Expected key length %d, got %d", KeySize, len(key))
	}
}

func TestDeriveKeyPBKDF2(t *testing.T) {
	config := DefaultConfig()
	config.KDF = KDFPBKDF2
	config.Iterations = PBKDF2Iter

	m := NewManager(config)

	salt := make([]byte, SaltSize)
	key := m.DeriveKey("password123", salt)

	if len(key) != KeySize {
		t.Errorf("Expected key length %d, got %d", KeySize, len(key))
	}
}

func TestDeriveKeyConsistency(t *testing.T) {
	config := DefaultConfig()
	m := NewManager(config)

	salt := make([]byte, SaltSize)
	password := "test-password"

	key1 := m.DeriveKey(password, salt)
	key2 := m.DeriveKey(password, salt)

	if !bytes.Equal(key1, key2) {
		t.Error("Same password and salt should produce same key")
	}
}

func TestDeriveKeyDifferentPasswords(t *testing.T) {
	config := DefaultConfig()
	m := NewManager(config)

	salt := make([]byte, SaltSize)

	key1 := m.DeriveKey("password1", salt)
	key2 := m.DeriveKey("password2", salt)

	if bytes.Equal(key1, key2) {
		t.Error("Different passwords should produce different keys")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	config := DefaultConfig()
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	m, err := NewManagerWithKey(config, key, "")
	if err != nil {
		t.Fatalf("NewManagerWithKey failed: %v", err)
	}

	plaintext := []byte("Hello, World! This is a secret message.")

	encData, err := m.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encData.Algorithm != AlgorithmAES256GCM {
		t.Errorf("Expected algorithm %s, got %s", AlgorithmAES256GCM, encData.Algorithm)
	}

	if len(encData.Nonce) != NonceSize {
		t.Errorf("Expected nonce size %d, got %d", NonceSize, len(encData.Nonce))
	}

	if len(encData.Ciphertext) == 0 {
		t.Error("Ciphertext should not be empty")
	}

	decrypted, err := m.Decrypt(encData)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Error("Decrypted text does not match original")
	}
}

func TestEncryptEmptyPlaintext(t *testing.T) {
	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	encData, err := m.Encrypt([]byte{})
	if err != nil {
		t.Fatalf("Encrypt failed for empty plaintext: %v", err)
	}

	decrypted, err := m.Decrypt(encData)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if len(decrypted) != 0 {
		t.Error("Expected empty decrypted text")
	}
}

func TestEncryptNoKey(t *testing.T) {
	config := DefaultConfig()
	m := NewManager(config)

	_, err := m.Encrypt([]byte("test"))
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestDecryptNoKey(t *testing.T) {
	config := DefaultConfig()
	m := NewManager(config)

	encData := &EncryptedData{
		Algorithm:  AlgorithmAES256GCM,
		Nonce:      make([]byte, NonceSize),
		Ciphertext: []byte("dummy"),
	}

	_, err := m.Decrypt(encData)
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestDecryptInvalidNonce(t *testing.T) {
	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	encData := &EncryptedData{
		Algorithm:  AlgorithmAES256GCM,
		Nonce:      make([]byte, 8),
		Ciphertext: []byte("dummy"),
	}

	_, err := m.Decrypt(encData)
	if err != ErrInvalidCiphertext {
		t.Errorf("Expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	plaintext := []byte("secret message")
	encData, _ := m.Encrypt(plaintext)

	encData.Ciphertext[0] ^= 0xFF

	_, err := m.Decrypt(encData)
	if err != ErrDecryptionFailed {
		t.Errorf("Expected ErrDecryptionFailed, got %v", err)
	}
}

func TestEncryptDecryptString(t *testing.T) {
	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	original := "This is a secret string value!"

	encrypted, err := m.EncryptString(original)
	if err != nil {
		t.Fatalf("EncryptString failed: %v", err)
	}

	if len(encrypted) == 0 {
		t.Error("Encrypted string should not be empty")
	}

	decrypted, err := m.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString failed: %v", err)
	}

	if original != decrypted {
		t.Errorf("Expected '%s', got '%s'", original, decrypted)
	}
}

func TestGenerateKey(t *testing.T) {
	key1, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if len(key1) != KeySize {
		t.Errorf("Expected key size %d, got %d", KeySize, len(key1))
	}

	key2, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if bytes.Equal(key1, key2) {
		t.Error("Two generated keys should not be equal")
	}
}

func TestEncryptWithPassword(t *testing.T) {
	config := DefaultConfig()
	plaintext := []byte("Secret data")

	encData, err := EncryptWithPassword(plaintext, "my-password", config)
	if err != nil {
		t.Fatalf("EncryptWithPassword failed: %v", err)
	}

	if encData.Algorithm != AlgorithmAES256GCM {
		t.Errorf("Expected algorithm %s", AlgorithmAES256GCM)
	}

	if len(encData.Salt) != SaltSize {
		t.Errorf("Expected salt size %d, got %d", SaltSize, len(encData.Salt))
	}
}

func TestDecryptWithPassword(t *testing.T) {
	config := DefaultConfig()
	plaintext := []byte("Secret data")
	password := "my-password"

	encData, err := EncryptWithPassword(plaintext, password, config)
	if err != nil {
		t.Fatalf("EncryptWithPassword failed: %v", err)
	}

	decrypted, err := DecryptWithPassword(encData, password, config)
	if err != nil {
		t.Fatalf("DecryptWithPassword failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Error("Decrypted text does not match original")
	}
}

func TestDecryptWithWrongPassword(t *testing.T) {
	config := DefaultConfig()
	plaintext := []byte("Secret data")

	encData, err := EncryptWithPassword(plaintext, "correct-password", config)
	if err != nil {
		t.Fatalf("EncryptWithPassword failed: %v", err)
	}

	_, err = DecryptWithPassword(encData, "wrong-password", config)
	if err == nil {
		t.Error("Expected decryption to fail with wrong password")
	}
}

func TestValidateKey(t *testing.T) {
	validKey := make([]byte, KeySize)
	invalidKey := make([]byte, 16)

	if !ValidateKey(validKey) {
		t.Error("32-byte key should be valid")
	}

	if ValidateKey(invalidKey) {
		t.Error("16-byte key should be invalid")
	}
}

func TestValidateEncryptedData(t *testing.T) {
	tests := []struct {
		name    string
		data    *EncryptedData
		wantErr bool
	}{
		{
			name: "valid data",
			data: &EncryptedData{
				Algorithm:  AlgorithmAES256GCM,
				Nonce:      make([]byte, NonceSize),
				Ciphertext: []byte("encrypted data"),
			},
			wantErr: false,
		},
		{
			name: "invalid algorithm",
			data: &EncryptedData{
				Algorithm:  "invalid",
				Nonce:      make([]byte, NonceSize),
				Ciphertext: []byte("encrypted data"),
			},
			wantErr: true,
		},
		{
			name: "invalid nonce size",
			data: &EncryptedData{
				Algorithm:  AlgorithmAES256GCM,
				Nonce:      make([]byte, 8),
				Ciphertext: []byte("encrypted data"),
			},
			wantErr: true,
		},
		{
			name: "empty ciphertext",
			data: &EncryptedData{
				Algorithm:  AlgorithmAES256GCM,
				Nonce:      make([]byte, NonceSize),
				Ciphertext: []byte{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEncryptedData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEncryptedData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptedDataJSON(t *testing.T) {
	original := &EncryptedData{
		Algorithm:  AlgorithmAES256GCM,
		KDF:        KDFArgon2,
		Salt:       []byte("salt123"),
		Nonce:      []byte("nonce123"),
		Ciphertext: []byte("ciphertext123"),
		KeyHint:    "test-key",
		Version:    1,
	}

	jsonStr, err := original.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	var restored EncryptedData
	if err := restored.FromJSON(jsonStr); err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if restored.Algorithm != original.Algorithm {
		t.Errorf("Algorithm mismatch")
	}

	if !bytes.Equal(restored.Salt, original.Salt) {
		t.Errorf("Salt mismatch")
	}
}

func TestEncryptedDataBase64(t *testing.T) {
	original := &EncryptedData{
		Algorithm:  AlgorithmAES256GCM,
		KDF:        KDFArgon2,
		Salt:       []byte("salt123"),
		Nonce:      []byte("nonce123"),
		Ciphertext: []byte("ciphertext123"),
		Version:    1,
	}

	b64, err := original.ToBase64()
	if err != nil {
		t.Fatalf("ToBase64 failed: %v", err)
	}

	var restored EncryptedData
	if err := restored.FromBase64(b64); err != nil {
		t.Fatalf("FromBase64 failed: %v", err)
	}

	if !bytes.Equal(restored.Salt, original.Salt) {
		t.Errorf("Salt mismatch after base64 roundtrip")
	}
}

func TestMemoryKeyStore(t *testing.T) {
	store := NewMemoryKeyStore()

	key := []byte("test-key-12345678901234567890")

	err := store.StoreKey("key1", key)
	if err != nil {
		t.Fatalf("StoreKey failed: %v", err)
	}

	retrieved, err := store.GetKey("key1")
	if err != nil {
		t.Fatalf("GetKey failed: %v", err)
	}

	if !bytes.Equal(key, retrieved) {
		t.Error("Retrieved key does not match original")
	}

	keys, err := store.ListKeys()
	if err != nil {
		t.Fatalf("ListKeys failed: %v", err)
	}

	if len(keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(keys))
	}

	err = store.DeleteKey("key1")
	if err != nil {
		t.Fatalf("DeleteKey failed: %v", err)
	}

	_, err = store.GetKey("key1")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestMemoryKeyStoreNotFound(t *testing.T) {
	store := NewMemoryKeyStore()

	_, err := store.GetKey("nonexistent")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestFileKeyStore(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileKeyStore(tmpDir)

	key := []byte("test-key-12345678901234567890")

	err := store.StoreKey("key1", key)
	if err != nil {
		t.Fatalf("StoreKey failed: %v", err)
	}

	retrieved, err := store.GetKey("key1")
	if err != nil {
		t.Fatalf("GetKey failed: %v", err)
	}

	if !bytes.Equal(key, retrieved) {
		t.Error("Retrieved key does not match original")
	}

	keys, err := store.ListKeys()
	if err != nil {
		t.Fatalf("ListKeys failed: %v", err)
	}

	if len(keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(keys))
	}

	err = store.DeleteKey("key1")
	if err != nil {
		t.Fatalf("DeleteKey failed: %v", err)
	}

	_, err = store.GetKey("key1")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestFileKeyStoreNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	store := NewFileKeyStore(tmpDir)

	_, err := store.GetKey("nonexistent")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestSecureConfig(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")

	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	sc := NewSecureConfig(m, tmpFile)

	err := sc.Set("api_key", "secret-api-key", true)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	err = sc.Set("debug", true, false)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if !sc.IsEncrypted("api_key") {
		t.Error("api_key should be encrypted")
	}

	if sc.IsEncrypted("debug") {
		t.Error("debug should not be encrypted")
	}

	value, err := sc.GetString("api_key")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}

	if value != "secret-api-key" {
		t.Errorf("Expected 'secret-api-key', got '%s'", value)
	}

	keys := sc.ListKeys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	sc.Delete("api_key")

	if sc.IsEncrypted("api_key") {
		t.Error("api_key should be deleted")
	}
}

func TestSecureConfigSaveLoad(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")

	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	sc := NewSecureConfig(m, tmpFile)

	sc.Set("api_key", "secret-key", true)
	sc.Set("host", "localhost", false)
	sc.Set("port", 8080, false)

	err := sc.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	sc2 := NewSecureConfig(m, tmpFile)
	err = sc2.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	apiKey, err := sc2.GetString("api_key")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}

	if apiKey != "secret-key" {
		t.Errorf("Expected 'secret-key', got '%s'", apiKey)
	}
}

func TestSecureConfigExportImport(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")

	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	sc := NewSecureConfig(m, tmpFile)

	sc.Set("api_key", "secret-key", true)
	sc.Set("debug", true, false)

	exported, err := sc.Export()
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	var exportData map[string]interface{}
	if err := json.Unmarshal(exported, &exportData); err != nil {
		t.Fatalf("Failed to unmarshal exported data: %v", err)
	}

	if exportData["api_key"] != "secret-key" {
		t.Errorf("Expected 'secret-key', got '%v'", exportData["api_key"])
	}

	sc2 := NewSecureConfig(m, filepath.Join(t.TempDir(), "config2.json"))
	err = sc2.Import(exported, true)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	value, err := sc2.GetString("api_key")
	if err != nil {
		t.Fatalf("GetString failed: %v", err)
	}

	if value != "secret-key" {
		t.Errorf("Expected 'secret-key', got '%s'", value)
	}
}

func TestEncryptDecryptFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "plaintext.txt")
	dstFile := filepath.Join(tmpDir, "encrypted.json")

	plaintext := []byte("This is a secret file content!")
	if err := os.WriteFile(srcFile, plaintext, 0644); err != nil {
		t.Fatalf("Failed to write source file: %v", err)
	}

	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	err := EncryptFile(srcFile, dstFile, m)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Error("Encrypted file was not created")
	}

	decrypted, err := DecryptFile(dstFile, m)
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Error("Decrypted content does not match original")
	}
}

func TestSecureConfigGetNotFound(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")

	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	sc := NewSecureConfig(m, tmpFile)

	_, err := sc.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent key")
	}
}

func TestSecureConfigGetStringNotString(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "config.json")

	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	sc := NewSecureConfig(m, tmpFile)
	sc.Set("number", 12345, false)

	_, err := sc.GetString("number")
	if err == nil {
		t.Error("Expected error when getting non-string as string")
	}
}

func BenchmarkEncrypt(b *testing.B) {
	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")

	plaintext := []byte("This is a benchmark test string for encryption performance.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encData, _ := m.Encrypt(plaintext)
		_, _ = m.Decrypt(encData)
	}
}

func BenchmarkDeriveKey(b *testing.B) {
	config := DefaultConfig()
	m := NewManager(config)
	salt := make([]byte, SaltSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.DeriveKey("benchmark-password", salt)
	}
}

func BenchmarkSecureConfigSet(b *testing.B) {
	tmpFile := filepath.Join(b.TempDir(), "config.json")
	config := DefaultConfig()
	key, _ := GenerateKey()
	m, _ := NewManagerWithKey(config, key, "")
	sc := NewSecureConfig(m, tmpFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sc.Set("key", "value", true)
	}
}
