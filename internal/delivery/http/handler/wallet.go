package handler

import (
	"mpc/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletUseCase *usecase.WalletUseCase
}

func NewWalletHandler(walletUseCase *usecase.WalletUseCase) *WalletHandler {
	return &WalletHandler{walletUseCase: walletUseCase}
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	// Implement the handler logic here
	c.JSON(http.StatusOK, gin.H{"message": "Wallet created"})
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	// Implement the handler logic here
	c.JSON(http.StatusOK, gin.H{"message": "Wallet retrieved"})
}