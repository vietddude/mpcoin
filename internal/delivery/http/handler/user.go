package handler

import (
	"mpc/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	// Implement the handler logic here
	c.JSON(http.StatusOK, gin.H{"message": "User retrieved"})
}
