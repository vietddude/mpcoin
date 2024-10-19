package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mpc/internal/infrastructure/redis"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	config      *JWTConfig
	signingKey  []byte
	redisClient redis.RedisClient
}

func NewJWTService(config *JWTConfig, redisClient redis.RedisClient) *JWTService {
	return &JWTService{
		config:      config,
		signingKey:  []byte(config.SecretKey),
		redisClient: redisClient,
	}
}

func (s *JWTService) generateToken(ctx context.Context, userID uuid.UUID, tokenType TokenType) (string, error) {
	var duration time.Duration
	if tokenType == AccessToken {
		duration = s.config.AccessTokenDuration
	} else {
		duration = s.config.RefreshTokenDuration
	}

	now := time.Now()
	claims := &JWTClaims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", err
	}
	// convert userID to string
	userIDString := userID.String()
	// Store token in Redis
	key := fmt.Sprintf("%s:%s", tokenType, userIDString)
	err = s.redisClient.Set(ctx, key, tokenString, duration)
	if err != nil {
		return "", fmt.Errorf("failed to store token in Redis: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) GenerateAccessToken(ctx context.Context, userID uuid.UUID) (string, error) {
	return s.generateToken(ctx, userID, AccessToken)
}

func (s *JWTService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	return s.generateToken(ctx, userID, RefreshToken)
}

func (s *JWTService) ValidateToken(ctx context.Context, tokenString string, tokenType TokenType) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	fmt.Println("validate token...")
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		if claims.Type != tokenType {
			return nil, errors.New("invalid token type")
		}

		// Validate token against Redis
		key := fmt.Sprintf("%s:%s", tokenType, claims.UserID.String())
		storedToken, err := s.redisClient.Get(ctx, key)
		if err != nil {
			if err.Error() == "redis: nil" {
				return nil, errors.New("token not found in Redis")
			}
			return nil, fmt.Errorf("failed to retrieve token from Redis: %w", err)
		}

		if storedToken != tokenString {
			return nil, errors.New("invalid token")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.ValidateToken(ctx, refreshToken, RefreshToken)
	if err != nil {
		return "", "", err
	}
	fmt.Println("claims", claims)

	// Generate new access and refresh tokens
	newAccessToken, err := s.GenerateAccessToken(ctx, claims.UserID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.GenerateRefreshToken(ctx, claims.UserID)
	if err != nil {
		return "", "", err
	}

	// Invalidate old refresh token
	s.InvalidateToken(ctx, claims.UserID, RefreshToken)

	return newAccessToken, newRefreshToken, nil
}

func (s *JWTService) InvalidateToken(ctx context.Context, userID uuid.UUID, tokenType TokenType) error {
	key := fmt.Sprintf("%s:%s", tokenType, userID.String())
	err := s.redisClient.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to invalidate token in Redis: %w", err)
	}
	return nil
}
