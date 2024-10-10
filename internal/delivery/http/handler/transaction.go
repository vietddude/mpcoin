package handler

import (
	"mpc/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUseCase *usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{transactionUseCase: transactionUseCase}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	// Implement the handler logic here
	c.JSON(http.StatusOK, gin.H{"message": "Transaction created"})
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	// Implement the handler logic here
	c.JSON(http.StatusOK, gin.H{"message": "Transaction retrieved"})
}
