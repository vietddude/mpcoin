package domain

import "time"

type Balance struct {
	ID        int64
	WalletID  int64
	ChainID   int64
	TokenID   int64
	Balance   float64
	UpdatedAt time.Time
}
