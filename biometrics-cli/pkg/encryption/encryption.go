package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

var (
	ErrInvalidKey        = errors.New("invalid encryption key")
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	ErrDecryptionFailed  = errors.New("decryption failed")
	ErrKeyNotFound       = errors.New("encryption key not found")
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrAlreadyLocked     = errors.New("config already locked")
	ErrNotLocked         = errors.New("config not locked")
)

const (
	KeySize       = 32
	NonceSize     = 12
	SaltSize      = 16
	Argon2Time    = 3
	Argon2Memory  = 64 * 1024
	Argon2Threads = 4
	PBKDF2Iter    = 100000
)

type EncryptionAlgorithm string

const (
	AlgorithmAES256GCM EncryptionAlgorithm = "AES-256-GCM"
)

type KDFType string

const (
	KDFArgon2 KDFType = "argon2"
	KDFPBKDF2 KDFType = "pbkdf2"
)

type Config struct {
	Algorithm  EncryptionAlgorithm `json:"algorithm"`
	KDF        KDFType             `json:"kdf"`
	KeySize    int                 `json:"key_size"`
	NonceSize  int                 `json:"nonce_size"`
	SaltSize   int                 `json:"salt_size"`
	Iterations int                 `json:"iterations,omitempty"`
	Memory     uint32              `json:"memory,omitempty"`
	Time       uint32              `json:"time,omitempty"`
	Threads    uint8               `json:"threads,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		Algorithm: AlgorithmAES256GCM,
		KDF:       KDFArgon2,
		KeySize:   KeySize,
		NonceSize: NonceSize,
		SaltSize:  SaltSize,
		Time:      Argon2Time,
		Memory:    Argon2Memory,
		Threads:   Argon2Threads,
	}
}

type EncryptedData struct {
	Algorithm  EncryptionAlgorithm `json:"algorithm"`
	KDF        KDFType             `json:"kdf"`
	Salt       []byte              `json:"salt"`
	Nonce      []byte              `json:"nonce"`
	Ciphertext []byte              `json:"ciphertext"`
	KeyHint    string              `json:"key_hint,omitempty"`
	Version    int                 `json:"version"`
}

type EncryptedConfig struct {
	Data      map[string]EncryptedData `json:"data"`
	Metadata  map[string]string        `json:"metadata,omitempty"`
	Version   int                      `json:"version"`
	CreatedAt string                   `json:"created_at,omitempty"`
	UpdatedAt string                   `json:"updated_at,omitempty"`
}

type Manager struct {
	config    Config
	masterKey []byte
	keyHint   string
	mu        sync.RWMutex
}

func NewManager(config Config) *Manager {
	return &Manager{
		config: config,
	}
}

func NewManagerWithKey(config Config, masterKey []byte, keyHint string) (*Manager, error) {
	if len(masterKey) != KeySize {
		return nil, ErrInvalidKey
	}

	m := &Manager{
		config:    config,
		masterKey: make([]byte, len(masterKey)),
		keyHint:   keyHint,
	}
	copy(m.masterKey, masterKey)

	return m, nil
}

func (m *Manager) SetMasterKey(key []byte, keyHint string) error {
	if len(key) != KeySize {
		return ErrInvalidKey
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.masterKey = make([]byte, len(key))
	copy(m.masterKey, key)
	m.keyHint = keyHint

	return nil
}

func (m *Manager) DeriveKey(password string, salt []byte) []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch m.config.KDF {
	case KDFArgon2:
		return argon2.IDKey(
			[]byte(password),
			salt,
			m.config.Time,
			m.config.Memory,
			m.config.Threads,
			uint32(m.config.KeySize),
		)
	case KDFPBKDF2:
		if m.config.Iterations == 0 {
			m.config.Iterations = PBKDF2Iter
		}
		return pbkdf2.Key(
			[]byte(password),
			salt,
			m.config.Iterations,
			m.config.KeySize,
			sha256.New,
		)
	default:
		return argon2.IDKey(
			[]byte(password),
			salt,
			Argon2Time,
			Argon2Memory,
			Argon2Threads,
			KeySize,
		)
	}
}

func (m *Manager) Encrypt(plaintext []byte) (*EncryptedData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.masterKey) == 0 {
		return nil, ErrKeyNotFound
	}

	salt := make([]byte, m.config.SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	nonce := make([]byte, m.config.NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	return &EncryptedData{
		Algorithm:  m.config.Algorithm,
		KDF:        m.config.KDF,
		Salt:       salt,
		Nonce:      nonce,
		Ciphertext: ciphertext,
		KeyHint:    m.keyHint,
		Version:    1,
	}, nil
}

func (m *Manager) Decrypt(encData *EncryptedData) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.masterKey) == 0 {
		return nil, ErrKeyNotFound
	}

	if encData.Algorithm != AlgorithmAES256GCM {
		return nil, fmt.Errorf("unsupported algorithm: %s", encData.Algorithm)
	}

	if len(encData.Nonce) != m.config.NonceSize {
		return nil, ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, encData.Nonce, encData.Ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

func (m *Manager) EncryptString(plaintext string) (string, error) {
	encData, err := m.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(encData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal encrypted data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(jsonData), nil
}

func (m *Manager) DecryptString(ciphertext string) (string, error) {
	jsonData, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	var encData EncryptedData
	if err := json.Unmarshal(jsonData, &encData); err != nil {
		return "", fmt.Errorf("failed to unmarshal encrypted data: %w", err)
	}

	plaintext, err := m.Decrypt(&encData)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

type SecureConfig struct {
	manager   *Manager
	filePath  string
	data      map[string]interface{}
	encrypted map[string]bool
	mu        sync.RWMutex
}

func NewSecureConfig(manager *Manager, filePath string) *SecureConfig {
	return &SecureConfig{
		manager:   manager,
		filePath:  filePath,
		data:      make(map[string]interface{}),
		encrypted: make(map[string]bool),
	}
}

func (sc *SecureConfig) Set(key string, value interface{}, encrypt bool) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if encrypt {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}

		encData, err := sc.manager.Encrypt(jsonValue)
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}

		sc.data[key] = encData
		sc.encrypted[key] = true
	} else {
		sc.data[key] = value
		sc.encrypted[key] = false
	}

	return nil
}

func (sc *SecureConfig) Get(key string) (interface{}, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	value, exists := sc.data[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	if sc.encrypted[key] {
		encData, ok := value.(*EncryptedData)
		if !ok {
			return nil, ErrInvalidConfig
		}

		plaintext, err := sc.manager.Decrypt(encData)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt value: %w", err)
		}

		var result interface{}
		if err := json.Unmarshal(plaintext, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal value: %w", err)
		}

		return result, nil
	}

	return value, nil
}

func (sc *SecureConfig) GetString(key string) (string, error) {
	value, err := sc.Get(key)
	if err != nil {
		return "", err
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("value is not a string: %T", value)
	}

	return str, nil
}

func (sc *SecureConfig) Delete(key string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	delete(sc.data, key)
	delete(sc.encrypted, key)
}

func (sc *SecureConfig) ListKeys() []string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	keys := make([]string, 0, len(sc.data))
	for key := range sc.data {
		keys = append(keys, key)
	}
	return keys
}

func (sc *SecureConfig) IsEncrypted(key string) bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.encrypted[key]
}

func (sc *SecureConfig) Save() error {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	dir := filepath.Dir(sc.filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	encConfig := &EncryptedConfig{
		Data:    make(map[string]EncryptedData),
		Version: 1,
	}

	for key, value := range sc.data {
		if sc.encrypted[key] {
			encData, ok := value.(*EncryptedData)
			if ok {
				encConfig.Data[key] = *encData
			}
		} else {
			jsonValue, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal %s: %w", key, err)
			}

			encData, err := sc.manager.Encrypt(jsonValue)
			if err != nil {
				return fmt.Errorf("failed to encrypt %s: %w", key, err)
			}

			encConfig.Data[key] = *encData
		}
	}

	jsonData, err := json.MarshalIndent(encConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(sc.filePath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (sc *SecureConfig) Load() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	jsonData, err := os.ReadFile(sc.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	var encConfig EncryptedConfig
	if err := json.Unmarshal(jsonData, &encConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	sc.data = make(map[string]interface{})
	sc.encrypted = make(map[string]bool)

	for key, encData := range encConfig.Data {
		sc.data[key] = &encData
		sc.encrypted[key] = true
	}

	return nil
}

func (sc *SecureConfig) Export() ([]byte, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	export := make(map[string]interface{})

	for key := range sc.data {
		value, err := sc.Get(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s: %w", key, err)
		}
		export[key] = value
	}

	return json.MarshalIndent(export, "", "  ")
}

func (sc *SecureConfig) Import(data []byte, encryptAll bool) error {
	var values map[string]interface{}
	if err := json.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	for key, value := range values {
		if err := sc.Set(key, value, encryptAll); err != nil {
			return fmt.Errorf("failed to set %s: %w", key, err)
		}
	}

	return nil
}

func GenerateKey() ([]byte, error) {
	key := make([]byte, KeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

func GenerateKeyFromPassword(password string, config Config) ([]byte, error) {
	salt := make([]byte, config.SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	m := NewManager(config)
	return m.DeriveKey(password, salt), nil
}

func EncryptWithPassword(plaintext []byte, password string, config Config) (*EncryptedData, error) {
	salt := make([]byte, config.SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	m := NewManager(config)
	key := m.DeriveKey(password, salt)

	master, err := NewManagerWithKey(config, key, "")
	if err != nil {
		return nil, err
	}

	encData, err := master.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}

	encData.Salt = salt
	return encData, nil
}

func DecryptWithPassword(encData *EncryptedData, password string, config Config) ([]byte, error) {
	m := NewManager(config)
	key := m.DeriveKey(password, encData.Salt)

	master, err := NewManagerWithKey(config, key, "")
	if err != nil {
		return nil, err
	}

	return master.Decrypt(encData)
}

func EncryptFile(srcPath, dstPath string, manager *Manager) error {
	plaintext, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	encData, err := manager.Encrypt(plaintext)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %w", err)
	}

	jsonData, err := json.MarshalIndent(encData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted data: %w", err)
	}

	if err := os.WriteFile(dstPath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	return nil
}

func DecryptFile(srcPath string, manager *Manager) ([]byte, error) {
	jsonData, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted file: %w", err)
	}

	var encData EncryptedData
	if err := json.Unmarshal(jsonData, &encData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted data: %w", err)
	}

	return manager.Decrypt(&encData)
}

func ValidateKey(key []byte) bool {
	return len(key) == KeySize
}

func ValidateEncryptedData(data *EncryptedData) error {
	if data.Algorithm != AlgorithmAES256GCM {
		return fmt.Errorf("unsupported algorithm: %s", data.Algorithm)
	}

	if len(data.Nonce) != NonceSize {
		return fmt.Errorf("invalid nonce size: expected %d, got %d", NonceSize, len(data.Nonce))
	}

	if len(data.Ciphertext) == 0 {
		return ErrInvalidCiphertext
	}

	return nil
}

func (ed *EncryptedData) ToJSON() (string, error) {
	jsonData, err := json.Marshal(ed)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (ed *EncryptedData) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), ed)
}

func (ed *EncryptedData) ToBase64() (string, error) {
	jsonData, err := json.Marshal(ed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}

func (ed *EncryptedData) FromBase64(b64 string) error {
	jsonData, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, ed)
}

type KeyStore interface {
	StoreKey(keyID string, key []byte) error
	GetKey(keyID string) ([]byte, error)
	DeleteKey(keyID string) error
	ListKeys() ([]string, error)
}

type MemoryKeyStore struct {
	keys map[string][]byte
	mu   sync.RWMutex
}

func NewMemoryKeyStore() *MemoryKeyStore {
	return &MemoryKeyStore{
		keys: make(map[string][]byte),
	}
}

func (m *MemoryKeyStore) StoreKey(keyID string, key []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.keys[keyID] = make([]byte, len(key))
	copy(m.keys[keyID], key)
	return nil
}

func (m *MemoryKeyStore) GetKey(keyID string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key, exists := m.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	result := make([]byte, len(key))
	copy(result, key)
	return result, nil
}

func (m *MemoryKeyStore) DeleteKey(keyID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.keys, keyID)
	return nil
}

func (m *MemoryKeyStore) ListKeys() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0, len(m.keys))
	for keyID := range m.keys {
		keys = append(keys, keyID)
	}
	return keys, nil
}

type FileKeyStore struct {
	basePath string
	mu       sync.RWMutex
}

func NewFileKeyStore(basePath string) *FileKeyStore {
	return &FileKeyStore{basePath: basePath}
}

func (f *FileKeyStore) StoreKey(keyID string, key []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := os.MkdirAll(f.basePath, 0700); err != nil {
		return err
	}

	keyPath := filepath.Join(f.basePath, keyID+".key")

	encrypted, err := encryptKeyFile(key)
	if err != nil {
		return err
	}

	return os.WriteFile(keyPath, encrypted, 0600)
}

func (f *FileKeyStore) GetKey(keyID string) ([]byte, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	keyPath := filepath.Join(f.basePath, keyID+".key")

	data, err := os.ReadFile(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	return decryptKeyFile(data)
}

func (f *FileKeyStore) DeleteKey(keyID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	keyPath := filepath.Join(f.basePath, keyID+".key")
	return os.Remove(keyPath)
}

func (f *FileKeyStore) ListKeys() ([]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	entries, err := os.ReadDir(f.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	keys := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".key" {
			keys = append(keys, entry.Name()[:len(entry.Name())-4])
		}
	}
	return keys, nil
}

func encryptKeyFile(key []byte) ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	machineID := getMachineIdentifier()
	derivedKey := pbkdf2.Key(machineID, salt, PBKDF2Iter, KeySize, sha256.New)

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, key, nil)

	result := make([]byte, 0, SaltSize+NonceSize+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

func decryptKeyFile(data []byte) ([]byte, error) {
	if len(data) < SaltSize+NonceSize {
		return nil, ErrInvalidCiphertext
	}

	salt := data[:SaltSize]
	nonce := data[SaltSize : SaltSize+NonceSize]
	ciphertext := data[SaltSize+NonceSize:]

	machineID := getMachineIdentifier()
	derivedKey := pbkdf2.Key(machineID, salt, PBKDF2Iter, KeySize, sha256.New)

	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func getMachineIdentifier() []byte {
	hostname, _ := os.Hostname()
	combined := hostname + "-biometrics-encryption-key"
	return []byte(combined)
}
