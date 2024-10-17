package domain

import "time"

type Chain struct {
	ID             int64
	Name           string
	ChainID        string
	RPCURL         string
	NativeCurrency string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
