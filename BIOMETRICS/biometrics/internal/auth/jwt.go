package auth

import (
	"biometrics/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret        string
	tokenExpire   time.Duration
	refreshExpire time.Duration
	logger        *utils.Logger
}

type Claims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func NewJWTManager(secret string, tokenExpire, refreshExpire time.Duration) *JWTManager {
	return &JWTManager{
		secret:        secret,
		tokenExpire:   tokenExpire,
		refreshExpire: refreshExpire,
		logger:        utils.NewLogger("info", "development"),
	}
}

func (m *JWTManager) GenerateTokenPair() (*TokenPair, error) {
	return nil, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	return nil, nil
}

func (m *JWTManager) RefreshAccessToken(refreshToken string) (string, error) {
	return "", nil
}
