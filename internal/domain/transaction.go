package domain

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Transaction struct {
	ID           int64
	FromWalletID int64
	ToWalletID   int64
	Amount       pgtype.Numeric
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateTransactionParams struct {
	FromWalletID int64
	ToWalletID   int64
	Amount       pgtype.Numeric
	Status       string
}
