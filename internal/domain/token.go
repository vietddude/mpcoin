package domain

import "time"

type Token struct {
	ID              int64
	ChainID         int64
	ContractAddress string
	Name            string
	Symbol          string
	Decimals        int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
