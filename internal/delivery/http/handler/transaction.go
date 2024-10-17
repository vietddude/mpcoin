package handler

import (
	"mpc/internal/domain"
	"mpc/internal/usecase"
	"mpc/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TxnHandler struct {
	txnUC usecase.TxnUseCase
}

func NewTxnHandler(txnUC usecase.TxnUseCase) *TxnHandler {
	return &TxnHandler{txnUC: txnUC}
}

func (h *TxnHandler) GetTransactions(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Transactions retrieved"})
}

func (h *TxnHandler) CreateTransaction(c *gin.Context) {
	// Get userID from auth middleware
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse request body
	var req domain.CreateTxnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Set userID from authenticated user
	parsedUserID, err := uuid.Parse(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	// Create transaction
	txnID, err := h.txnUC.CreateTransaction(c.Request.Context(), parsedUserID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{
		"message": "Transaction created successfully",
		"txn_id":  txnID,
	})
}

func (h *TxnHandler) SubmitTransaction(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Transaction submitted"})
}
