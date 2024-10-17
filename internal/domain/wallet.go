package domain

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Address             string
	EncryptedPrivateKey []byte
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CreateWalletParams struct {
	UserID              uuid.UUID
	Address             string
	EncryptedPrivateKey []byte
}

type CreateWalletResponse struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	Address string    `json:"address"`
}
