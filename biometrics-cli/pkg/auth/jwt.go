package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("token expired")
	ErrInvalidPassword = errors.New("invalid password")
)

type JWTClaims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	TenantID  string   `json:"tenant_id"`
	SessionID string   `json:"session_id"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type UserManager interface {
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	ValidatePassword(email, password string) (*User, error)
}

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Roles        []string  `json:"roles"`
	TenantID     string    `json:"tenant_id"`
	MfaEnabled   bool      `json:"mfa_enabled"`
	MfaSecret    string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
}

type JWTManager struct {
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	accessTokenTTL      time.Duration
	refreshTokenTTL     time.Duration
	mu                  sync.RWMutex
	keyRotatedAt        time.Time
	keyRotationInterval time.Duration
}

func NewJWTManager(accessTokenTTL, refreshTokenTTL time.Duration) (*JWTManager, error) {
	privateKey, err := loadOrGenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to load key: %w", err)
	}

	return &JWTManager{
		privateKey:          privateKey,
		publicKey:           &privateKey.PublicKey,
		accessTokenTTL:      accessTokenTTL,
		refreshTokenTTL:     refreshTokenTTL,
		keyRotationInterval: 24 * time.Hour,
	}, nil
}

func loadOrGenerateKey() (*rsa.PrivateKey, error) {
	keyFile := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if keyFile == "" {
		keyFile = "/tmp/biometrics_jwt_key.pem"
	}

	if data, err := os.ReadFile(keyFile); err == nil {
		block, _ := pem.Decode(data)
		if block == nil {
			return nil, errors.New("failed to parse PEM block")
		}
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}

	privatePEM := pem.EncodeToMemory(&privBlock)
	if err := os.WriteFile(keyFile, privatePEM, 0600); err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (m *JWTManager) GenerateTokenPair(user *User, sessionID string) (*TokenPair, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	accessExpires := now.Add(m.accessTokenTTL)
	refreshExpires := now.Add(m.refreshTokenTTL)

	accessClaims := &JWTClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Roles:     user.Roles,
		TenantID:  user.TenantID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpires),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "biometrics-auth",
			Subject:   user.ID,
			ID:        sessionID + "_access",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(m.privateKey)
	if err != nil {
		return nil, err
	}

	refreshClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshExpires),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "biometrics-auth",
		Subject:   user.ID,
		ID:        sessionID + "_refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(m.privateKey)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpires,
	}, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.publicKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

func (m *JWTManager) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := m.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       claims.UserID,
		Email:    claims.Email,
		Roles:    claims.Roles,
		TenantID: claims.TenantID,
	}

	return m.GenerateTokenPair(user, claims.SessionID)
}

func (m *JWTManager) RotateKeys() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if time.Since(m.keyRotatedAt) < m.keyRotationInterval {
		return nil
	}

	newPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	m.privateKey = newPrivateKey
	m.publicKey = &newPrivateKey.PublicKey
	m.keyRotatedAt = time.Now()

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func ComputeFingerprint(data []byte) string {
	hash := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}
