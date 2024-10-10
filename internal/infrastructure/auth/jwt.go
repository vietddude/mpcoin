package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

type JWTClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTConfig(secretKey string, tokenDuration time.Duration) *JWTConfig {
	return &JWTConfig{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration,
	}
}
