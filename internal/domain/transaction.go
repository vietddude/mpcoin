package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusSuccess   Status = "success"
	StatusFailed    Status = "failed"
	StatusSubmitted Status = "submitted"
)

type Transaction struct {
	ID          uuid.UUID
	WalletID    uuid.UUID
	ChainID     uuid.UUID
	FromAddress string
	ToAddress   string
	Amount      string
	TokenID     uuid.UUID
	GasPrice    pgtype.Numeric
	GasLimit    int64
	Nonce       int64
	Status      Status
	TxHash      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateTransactionParams struct {
	ID          uuid.UUID
	WalletID    uuid.UUID
	ChainID     uuid.UUID
	FromAddress string
	ToAddress   string
	Amount      string
	TokenID     uuid.UUID
	GasPrice    pgtype.Numeric
	GasLimit    int64
	Nonce       int64
	UnsignedTx  string
	Status      Status
}

type SubmitTransactionParams struct {
	ID uuid.UUID
}

type CreateTxnRequest struct {
	WalletID  uuid.UUID `json:"wallet_id" binding:"required"`
	ChainID   uuid.UUID `json:"chain_id" binding:"required"`
	ToAddress string    `json:"to_address" binding:"required"`
	Amount    string    `json:"amount" binding:"required"`
	TokenID   uuid.UUID `json:"token_id" binding:"required"`
}

type CreateTxnResponse struct {
	ID uuid.UUID `json:"id"`
}

type SubmitTxnRequest struct {
	ID uuid.UUID `json:"txn_id" binding:"required"`
}

type SubmitTxnResponse struct {
	TxHash string `json:"tx_hash"`
}
