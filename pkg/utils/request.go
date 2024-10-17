// pkg/utils/utils.go

package utils

import (
	"errors"
	"strings"

	"mpc/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseRequest[T any](c *gin.Context) (T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, err
	}
	return req, nil
}

func GetAuthToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}
	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	tokenString := tokenParts[1]

	return tokenString, nil
}

func GetUserIDFromAuthToken(c *gin.Context, jwtService auth.JWTService) (uuid.UUID, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return uuid.Nil, errors.New("authorization header missing")
	}

	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return uuid.Nil, errors.New("invalid authorization header format")
	}

	tokenString := tokenParts[1]
	claims, err := jwtService.ValidateToken(c.Request.Context(), tokenString, auth.AccessToken)
	if err != nil {
		return uuid.Nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claims.UserID.String())
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID in token")
	}

	return userID, nil
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

func SuccessResponse(c *gin.Context, statusCode int, payload any) {
	c.JSON(statusCode, gin.H{"payload": payload})
}
