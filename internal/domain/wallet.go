package domain

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Wallet struct {
	ID                  int64
	UserID              int64
	Address             string
	PublicKey           string
	EncryptedPrivateKey string
	Balance             pgtype.Numeric
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CreateWalletParams struct {
	UserID              int64
	Address             string
	PublicKey           string
	EncryptedPrivateKey string
}
