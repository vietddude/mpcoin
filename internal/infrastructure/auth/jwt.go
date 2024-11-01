package auth

import (
	"mpc/internal/infrastructure/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

func NewJWTConfig(cfg *config.JWTConfig) *JWTConfig {
	return &JWTConfig{
		SecretKey:            cfg.SecretKey,
		AccessTokenDuration:  cfg.TokenDuration,
		RefreshTokenDuration: cfg.TokenDuration * 30,
	}
}
